window.onload = function () {

    const userId = crypto.randomUUID();
    console.log(userId);

    // WebSocket connection to localhost
    const ws = new WebSocket('ws://localhost:8080/ws');

    // Chat elements
    const chatContainer = document.getElementById('chat-container');
    const messageInput = document.getElementById('message');
    const sendButton = document.getElementById('send');

    // Add a message to the chat
    function addMessage(text) {
        const messageElem = document.createElement('div');
        messageElem.textContent = text;
        messageElem.classList.add('message');
        chatContainer.appendChild(messageElem);
        chatContainer.scrollTop = chatContainer.scrollHeight;

        return messageElem;
    }

    // Handle WebSocket messages
    ws.onmessage = function (event) {
        const data = JSON.parse(event.data);

        if (data.message) {
            msg = document.getElementById(data.message_id);
            if (!msg) {
                addMessage(data.message, data.id);
            }
        }
    };

    // Post a message to the server
    function postMessage(text) {
        return fetch('http://localhost:8080/enqueue', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-User-ID': userId
            },
            body: JSON.stringify({ message: text })
        }).then(response => {
            if (!response.ok) {
                console.error('Failed to send message');
            }
            return response.json();
        });
    }

    // Send button click event
    sendButton.addEventListener('click', function () {
        const message = messageInput.value.trim();
        if (message) {
            const msgId = crypto.randomUUID();
            console.log(msgId);

            // Display message in the chat
            msg = addMessage(message, msgId, true);
            msg.classList.add('own-message');
            msg.classList.add('loading');

            // Send message to the server
            postMessage(message).then(resp => {
                msg.setAttribute("id", resp.message_id)
                msg.classList.remove('loading');
            })

            // Clear the input field
            messageInput.value = '';
        }
    });

    // Send message when pressing Enter
    messageInput.addEventListener('keypress', function (event) {
        if (event.key === 'Enter') {
            sendButton.click();
        }
    });
}