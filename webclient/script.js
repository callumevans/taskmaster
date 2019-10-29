var Taskmaster = (function (exports) {
    'use strict';

    var Task = /** @class */ (function () {
        function Task(id, attributes) {
            this.id = id;
            this.attributes = attributes;
        }
        Object.defineProperty(Task.prototype, "Id", {
            get: function () {
                return this.id;
            },
            enumerable: true,
            configurable: true
        });
        Object.defineProperty(Task.prototype, "Attributes", {
            get: function () {
                return this.attributes;
            },
            enumerable: true,
            configurable: true
        });
        Task.FromMessage = function (message) {
            return new Task(message.Data.Task.id, message.Data.Task.attributes);
        };
        return Task;
    }());

    var Message = /** @class */ (function () {
        function Message(messageType, data) {
            this.messageType = messageType;
            this.data = data;
        }
        Object.defineProperty(Message.prototype, "MessageType", {
            get: function () {
                return this.messageType;
            },
            enumerable: true,
            configurable: true
        });
        Object.defineProperty(Message.prototype, "Data", {
            get: function () {
                return this.data;
            },
            enumerable: true,
            configurable: true
        });
        Message.FromJSON = function (json) {
            var parsed = JSON.parse(json);
            return new Message(parsed.messageType, parsed.message);
        };
        return Message;
    }());

    var MessageType;
    (function (MessageType) {
        MessageType["TaskReservationCreated"] = "task.reservation_created";
        MessageType["TaskAccepted"] = "task.accepted";
        MessageType["TaskWorkflowTimeout"] = "task.workflow_timeout";
        MessageType["TaskStageTimeout"] = "task.stage_timeout";
    })(MessageType || (MessageType = {}));

    var MemoryCache = /** @class */ (function () {
        function MemoryCache() {
            this.cachedItems = new Array();
            this.onExpiredCallback = null;
        }
        MemoryCache.prototype.has = function (key) {
            return !!this.cachedItems.find(function (item) { return item.key === key; });
        };
        MemoryCache.prototype.delete = function (key) {
            var item = this.cachedItems.find(function (item) { return item.key === key; });
            clearTimeout(item.timeoutFunc);
            this.cachedItems = this.cachedItems.filter(function (item) { return item.key !== key; });
        };
        MemoryCache.prototype.set = function (key, value, ttl) {
            var _this = this;
            if (this.cachedItems.find(function (item) { return item.key; })) {
                this.delete(key);
            }
            var timeOut = setTimeout(function () {
                console.debug("MemoryCache > Expiring cached task " + key);
                _this.cachedItems = _this.cachedItems.filter(function (item) { return item.key !== key; });
                if (_this.onExpiredCallback) {
                    _this.onExpiredCallback(key, value);
                }
            }, ttl);
            this.cachedItems.push({ key: key, value: value, timeoutFunc: timeOut });
        };
        return MemoryCache;
    }());

    var Client = /** @class */ (function () {
        function Client(instanceUrl, workerId) {
            var _this = this;
            this.activeTasks = [];
            this.incomingTasks = [];
            this.stageTimeoutSyncTasks = new MemoryCache();
            this.activeTaskAddedCallback = null;
            this.activeTaskRemovedCallback = null;
            this.incomingTaskAddedCallback = null;
            this.incomingTaskRemovedCallback = null;
            this.getActiveTasks = function () {
                return Object.assign([], _this.activeTasks);
            };
            this.getIncomingTasks = function () {
                return Object.assign([], _this.incomingTasks);
            };
            this.messageReceivedHandler = function (messageData) {
                console.debug("Message Received Handler > Received messaged", messageData);
                var message = Message.FromJSON(messageData.data);
                switch (message.MessageType) {
                    case MessageType.TaskReservationCreated:
                        _this.reservationCreatedHandler(Task.FromMessage(message));
                        break;
                    case MessageType.TaskWorkflowTimeout:
                        _this.taskWorkflowTimeoutHandler(Task.FromMessage(message));
                        break;
                    case MessageType.TaskStageTimeout:
                        _this.taskStageTimeoutHandler(Task.FromMessage(message));
                        break;
                }
            };
            this.stageTimeoutSyncTaskExpiredHandler = function (taskId) {
                console.debug("Stage Timeout Handler", taskId);
                _this.removeIncomingTask(taskId);
            };
            this.taskWorkflowTimeoutHandler = function (task) {
                console.debug("Task Workflow Timeout Handler", task);
                _this.removeIncomingTask(task.Id);
                Client.tryInvokeCallback(_this.incomingTaskRemovedCallback, task);
            };
            this.reservationCreatedHandler = function (task) {
                console.debug("Reservation Created Handler", task);
                if (!_this.incomingTasks.find(function (t) { return t.Id === task.Id; })) {
                    console.debug("Reservation Created Handler > Adding incoming task", task);
                    _this.incomingTasks.push(task);
                    Client.tryInvokeCallback(_this.incomingTaskAddedCallback, task);
                }
                // If a stage timeout has occurred and we have a sync event waiting just remove it
                if (_this.stageTimeoutSyncTasks.has(task.Id)) {
                    console.debug("Reservation Created Handler > Deleting stage timeout sync entry", task);
                    _this.stageTimeoutSyncTasks.delete(task.Id);
                }
            };
            this.taskStageTimeoutHandler = function (task) {
                // When a stage timeout happens, we set a timeout before expiring the task locally
                // If we receive another reservation for the task, we don't expire it
                console.debug("Task Stage Timeout Handler > Queueing timeout sync record for " + task.Id, task);
                _this.stageTimeoutSyncTasks.set(task.Id, task, Client.STAGE_TIMEOUT_SYNC_TTL);
            };
            this.removeIncomingTask = function (taskId) {
                _this.incomingTasks = _this.incomingTasks.filter(function (t) { return t.Id === taskId; });
            };
            this.connection = new WebSocket("ws://" + instanceUrl + "/ws?workerId=" + workerId);
            this.connection.onclose = function (evt) {
                console.log('Connection closed', evt);
            };
            this.connection.onmessage = this.messageReceivedHandler;
            this.stageTimeoutSyncTasks.onExpiredCallback = function (key, value) { return _this.stageTimeoutSyncTaskExpiredHandler(key); };
        }
        Client.tryInvokeCallback = function (fn, value) {
            if (fn) {
                fn(value);
            }
        };
        Client.STAGE_TIMEOUT_SYNC_TTL = 4000;
        return Client;
    }());

    exports.Client = Client;

    return exports;

}({}));
