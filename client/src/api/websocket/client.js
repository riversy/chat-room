export const LoginRequestMessage = 1;
export const LoginResponseMessage = 2;
export const LogoutRequestMessage = 3;
export const LogoutResponseMessage = 4;
export const UserEnterMessage = 5;
export const UserQuitMessage = 6;
export const ClientTextMessage = 7;
export const ServerTextMessage = 8;
export const ParticipantsQtyUpdate = 9;

export const ConnectedEvent = 'connected';
export const DisconnectedEvent = 'disconnected';
export const ErrorEvent = 'error';

export class WebsocketClient {
    constructor(uri, handlers) {
        this.setHandlers(handlers);
        this.initSocket(uri)
    }

    setHandlers(handlers) {
        this.handlers = handlers;
    }

    initSocket(uri) {
        const socket = new WebSocket(uri);
        socket.onopen = () => {
            console.log("Connected");
            this.handleEvent(ConnectedEvent)
        };

        socket.onclose = (event) => {
            console.log("Disconnected: ", event);
            this.handleEvent(DisconnectedEvent, event)
        };

        socket.onerror = (error) => {
            console.log("Socket Error: ", error);
            this.handleEvent(ErrorEvent, error)
        };

        socket.onmessage = message => {
            console.log("Message arrived:", message);
            this.handleMessage(message);
        };

        this.socket = socket;
    }

    /**
     * @param Array props
     */
    handleEvent(...props) {
        const eventIdentifier = props.shift()
        const handler = this.handlers[eventIdentifier] || null;
        if (!handler) {
            return;
        }
        handler(props);
    }

    /**
     * @param event
     */
    handleMessage(event) {
        const {data = null} = event;
        if (!data) {
            return;
        }
        let message = {};
        try {
            message = JSON.parse(data);
            const {type, payload = null} = message;
            const handler = this.handlers[type] || null;
            if (!handler) {
                return;
            }
            handler(payload)
        } catch (err) {
            console.error(err);
        }
    }

    /**
     * @param messageType
     * @param message
     */
    sendMessage(messageType, message) {
        const transport = JSON.stringify({
            "type": messageType,
            "payload": message
        });
        this.socket.send(transport);
    }
}
