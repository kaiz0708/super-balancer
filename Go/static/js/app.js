let servers = [];

function addServer() {
    const url = document.getElementById('urlInput').value;
    const weight = document.getElementById('weightInput').value;
    const statusElement = document.getElementById('status');

    if (!url || !weight) {
        statusElement.textContent = 'Please fill in both URL and Weight.';
        statusElement.className = 'mt-4 text-center text-sm text-red-500';
        return;
    }

    const server = {
        url: url,
        healthy: true,
        avgLatency: 0,
        requestCount: 0,
        failureRate: 0,
        weight: parseInt(weight),
        activeConnections: 0
    };

    servers.push(server);
    updateServerList();
    
    // Clear inputs
    document.getElementById('urlInput').value = '';
    document.getElementById('weightInput').value = '';
    statusElement.textContent = 'Server added successfully!';
    statusElement.className = 'mt-4 text-center text-sm text-green-500';
}

function updateServerList() {
    const serverListItems = document.getElementById('serverListItems');
    serverListItems.innerHTML = '';

    servers.forEach(server => {
        const row = document.createElement('tr');
        row.className = 'border-t border-gray-200';
        row.innerHTML = `
            <td class="p-3">${server.url}</td>
            <td class="p-3">
                <span class="${server.healthy ? 'text-green-500' : 'text-red-500'} font-medium">
                    ${server.healthy ? 'Healthy' : 'Unhealthy'}
                </span>
            </td>
            <td class="p-3">${server.avgLatency}</td>
            <td class="p-3">${server.requestCount}</td>
            <td class="p-3">${server.failureRate}</td>
            <td class="p-3">${server.weight}</td>
            <td class="p-3">${server.activeConnections}</td>
        `;
        serverListItems.appendChild(row);
    });
}

function getServers() {
    const statusElement = document.getElementById('status');
    statusElement.textContent = 'Fetching servers...';
    statusElement.className = 'mt-4 text-center text-sm text-blue-500';

    fetch('http://localhost:8080/api/backends')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch servers');
            }
            return response.json();
        })
        .then(data => {
            servers = data.backends;
            updateServerList();
            statusElement.textContent = `Fetched ${servers.length} servers successfully!`;
            statusElement.className = 'mt-4 text-center text-sm text-green-500';
        })
        .catch(error => {
            statusElement.textContent = `Error: ${error.message}`;
            statusElement.className = 'mt-4 text-center text-sm text-red-500';
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
            servers = []; 
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