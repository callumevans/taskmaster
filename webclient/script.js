var Taskmaster = (function (exports) {
    'use strict';

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
            return new Message(json.messageType, json.message);
        };
        return Message;
    }());

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
        Task.FromJSON = function (json) {
            return new Task(json.id, json.attributes);
        };
        return Task;
    }());

    var Client = /** @class */ (function () {
        function Client(instanceUrl, workerId) {
            this.connection = new WebSocket("ws://" + instanceUrl + "/ws?workerId=" + workerId);
            this.connection.onclose = function (evt) {
                console.log('Connection closed', evt);
            };
            this.connection.onmessage = this.messageReceivedHandler;
        }
        Client.prototype.messageReceivedHandler = function (message) {
            var parsed = JSON.parse(message.data);
            var msg = Message.FromJSON(parsed);
            var task = Task.FromJSON(msg.Data);
            console.log('Task', task);
            console.log(msg);
        };
        return Client;
    }());

    exports.Client = Client;

    return exports;

}({}));
