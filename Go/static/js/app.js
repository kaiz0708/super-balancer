let servers = [];

function addServer() {
    const urlInput = document.getElementById('urlInput').value.trim();
    const weightInput = document.getElementById('weightInput').value.trim();
    const status = document.getElementById('status');

    // Validate inputs
    if (!urlInput || !weightInput) {
        status.textContent = 'Please enter both URL and weight.';
        status.className = 'text-red-500';
        return;
    }
    if (!urlInput.startsWith('http://') && !urlInput.startsWith('https://')) {
        status.textContent = 'URL must start with http:// or https://';
        status.className = 'text-red-500';
        return;
    }
    const weight = parseInt(weightInput);
    if (isNaN(weight) || weight <= 0) {
        status.textContent = 'Weight must be a positive number.';
        status.className = 'text-red-500';
        return;
    }

    // Add server to list
    servers.push({ url: urlInput, weight });
    updateServerList();

    // Clear inputs and show success
    document.getElementById('urlInput').value = '';
    document.getElementById('weightInput').value = '';
    status.textContent = 'Server added to list!';
    status.className = 'text-green-500';
}

function updateServerList() {
    const serverListItems = document.getElementById('serverListItems');
    serverListItems.innerHTML = '';
    servers.forEach((server, index) => {
        const li = document.createElement('li');
        li.textContent = `URL: ${server.url}, Weight: ${server.weight}`;
        serverListItems.appendChild(li);
    });
}

async function submitServers() {
    const status = document.getElementById('status');

    if (servers.length === 0) {
        status.textContent = 'No servers to submit.';
        status.className = 'text-red-500';
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/api/backends', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(servers),
        });

        const result = await response.json();
        if (response.ok) {
            status.textContent = result.message || 'Backends added successfully!';
            status.className = 'text-green-500';
            servers = []; // Clear the list
            updateServerList();
        } else {
            status.textContent = result.error || 'Failed to add backends.';
            status.className = 'text-red-500';
        }
    } catch (error) {
        status.textContent = 'Error: ' + error.message;
        status.className = 'text-red-500';
    }
}