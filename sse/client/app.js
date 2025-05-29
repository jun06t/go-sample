let eventSource;
const messageDiv = document.getElementById('message');
const startBtn = document.getElementById('startBtn');
const stopBtn = document.getElementById('stopBtn');
const statusDiv = document.getElementById('status');

let errorTimeout = null;

function updateStatus(message) {
    statusDiv.textContent = message;
}

function clearErrorTimeout() {
    if (errorTimeout) {
        clearTimeout(errorTimeout);
        errorTimeout = null;
    }
}

function startStream() {
    if (eventSource) eventSource.close();

    messageDiv.innerHTML = '';
    startBtn.disabled = true;
    stopBtn.disabled = false;

    eventSource = new EventSource('http://localhost:8080/stream');

    eventSource.onopen = clearErrorTimeout;

    eventSource.onmessage = function(event) {
        clearErrorTimeout();
        try {
            const data = JSON.parse(event.data);
            messageDiv.textContent = data.content;
            updateStatus(`Last received ID: ${data.id}`);
            if (data.done) {
                stopStream();
            }
        } catch (error) {
            console.error('Error parsing JSON:', error);
        }
    };

    eventSource.onerror = function(error) {
        console.error('Stream Error:', error);
        updateStatus('Connection error. Waiting 10 seconds before disconnect...');
        clearErrorTimeout();
        errorTimeout = setTimeout(() => {
            stopStream();
            updateStatus('Disconnected after error. Click Start to reconnect.');
        }, 10000);
    };
}

function stopStream() {
    if (eventSource) {
        eventSource.close();
        eventSource = null;
    }
    startBtn.disabled = false;
    stopBtn.disabled = true;
    clearErrorTimeout();
}

updateStatus('Ready to stream.'); 