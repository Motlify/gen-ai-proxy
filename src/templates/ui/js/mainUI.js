// src/templates/ui/js/mainUI.js

import { authenticatedFetch, handleDelete } from './api.js';

let apiKeyTableBody, modelTableBody, connectionTableBody, providerTableBody, conversationLogTableBody;
let navButtons, contentSections;

export const initMainUIElements = (apiKeysBody, modelsBody, connectionsBody, providersBody, conversationLogsBody, navBtns, contentSects) => {
    apiKeyTableBody = apiKeysBody;
    modelTableBody = modelsBody;
    connectionTableBody = connectionsBody;
    providerTableBody = providersBody;
    conversationLogTableBody = conversationLogsBody;
    navButtons = navBtns;
    contentSections = contentSects;
};

// --- UI Section Toggling Logic ---
export const showSection = (sectionId) => {
    if (contentSections.length > 0) {
        contentSections.forEach(section => {
            section.classList.add('hidden');
        });
        document.getElementById(sectionId).classList.remove('hidden');
    }
};

export const activateNavButton = (button) => {
    if (navButtons.length > 0) {
        navButtons.forEach(btn => {
            btn.classList.remove('bg-blue-600', 'text-white');
            btn.classList.add('bg-gray-200', 'dark:bg-gray-700', 'text-gray-800', 'dark:text-gray-200');
        });
        button.classList.add('bg-blue-600', 'text-white');
        button.classList.remove('bg-gray-200', 'dark:bg-gray-700', 'text-gray-800', 'dark:text-gray-200');
    }
};

// Fetch and display API Keys
export async function fetchApiKeys() {
    if (!apiKeyTableBody) return; // Ensure element exists
    apiKeyTableBody.innerHTML = '';
    try {
        const response = await authenticatedFetch('/api/api-keys');
        const result = await response.json();
        if (response.ok) {
            if (result.api_keys && result.api_keys.length > 0) {
                result.api_keys.forEach(key => {
                    const row = document.createElement('tr');
                    row.classList.add('table-row-hover', 'border-b', 'border-gray-200', 'dark:border-gray-700');
                    row.innerHTML = `
                        <td class="py-3 px-6 text-left">${key.id || 'N/A'}</td>
                        <td class="py-3 px-6 text-left">${key.name}</td>
                        <td class="py-3 px-6 text-left">${new Date(key.created_at).toLocaleString()}</td>
                        <td class="py-3 px-6 text-left">${key.last_used_at ? new Date(key.last_used_at).toLocaleString() : 'Never'}</td>
                        <td class="py-3 px-6 text-center">
                            <button class="delete-btn bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded focus:outline-none focus:shadow-outline" data-id="${key.id}" data-type="api-key">Delete</button>
                        </td>
                    `;
                    apiKeyTableBody.appendChild(row);
                });
                addDeleteEventListeners();
            } else {
                apiKeyTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">No API keys found.</td></tr>`;
            }
        } else {
            apiKeyTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">Error: ${result.error || 'Failed to fetch API keys'}</td></tr>`;
        }
    } catch (error) {
        console.error('Error fetching API keys:', error);
        apiKeyTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">An error occurred while fetching API keys.</td></tr>`;
    }
}

// Fetch and display Models
export async function fetchModels() {
    if (!modelTableBody) return; // Ensure element exists
    modelTableBody.innerHTML = '';
    try {
        const response = await authenticatedFetch('/api/models');
        const data = await response.json();
        if (response.ok) {
            if (data && data.length > 0) {
                data.forEach(model => {
                    const row = document.createElement('tr');
                    row.classList.add('table-row-hover', 'border-b', 'border-gray-200', 'dark:border-gray-700');
                    row.innerHTML = `
                        <td class="py-3 px-6 text-left">${model.id}</td>
                        <td class="py-3 px-6 text-left">${model.proxy_model_id}</td>
                        <td class="py-3 px-6 text-left">${model.provider_model_id}</td>
                        <td class="py-3 px-6 text-left">${model.connection_id}</td>
                        <td class="py-3 px-6 text-left">${model.price_input}</td>
                        <td class="py-3 px-6 text-left">${model.price_output}</td>
                        <td class="py-3 px-6 text-left">${model.thinking}</td>
                        <td class="py-3 px-6 text-left">${model.tools_usage}</td>
                        <td class="py-3 px-6 text-center">
                            <button class="delete-btn bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded focus:outline-none focus:shadow-outline" data-id="${model.id}" data-type="model">Delete</button>
                        </td>
                    `;
                    modelTableBody.appendChild(row);
                });
                addDeleteEventListeners();
            } else {
                modelTableBody.innerHTML = `<tr><td colspan="6" class="py-3 px-6 text-center">No models found.</td></tr>`;
            }
        } else {
            const errorData = await response.json();
            modelTableBody.innerHTML = `<tr><td colspan="6" class="py-3 px-6 text-center">Error: ${errorData.error || 'Failed to fetch models'}</td></tr>`;
        }
    } catch (error) {
        console.error('Error fetching models:', error);
        modelTableBody.innerHTML = `<tr><td colspan="6" class="py-3 px-6 text-center">An error occurred while fetching models.</td></tr>`;
    }
}

// Fetch and display Connections
export async function fetchConnections() {
    if (!connectionTableBody) return; // Ensure element exists
    connectionTableBody.innerHTML = '';
    try {
        const response = await authenticatedFetch('/api/connections');
        const data = await response.json();
        if (response.ok) {
            if (data.connections && data.connections.length > 0) {
                data.connections.forEach(conn => {
                    const row = document.createElement('tr');
                    row.classList.add('table-row-hover', 'border-b', 'border-gray-200', 'dark:border-gray-700');
                    row.innerHTML = `
                        <td class="py-3 px-6 text-left">${conn.id}</td>
                        <td class="py-3 px-6 text-left">${conn.name}</td>
                        <td class="py-3 px-6 text-left">${conn.provider}</td>
                        <td class="py-3 px-6 text-left">${new Date(conn.created_at).toLocaleString()}</td>
                        <td class="py-3 px-6 text-center">
                            <button class="delete-btn bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded focus:outline-none focus:shadow-outline" data-id="${conn.id}" data-type="connection">Delete</button>
                        </td>
                    `;
                    connectionTableBody.appendChild(row);
                });
                addDeleteEventListeners();
            } else {
                connectionTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">No connections found.</td></tr>`;
            }
        } else {
            const errorData = await response.json();
            connectionTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">Error: ${errorData.error || 'Failed to fetch connections'}</td></tr>`;
        }
    } catch (error) {
        console.error('Error fetching connections:', error);
        connectionTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">An error occurred while fetching connections.</td></tr>`;
    }
}

// Fetch and display Providers
export async function fetchProviders() {
    if (!providerTableBody) return; // Ensure element exists
    providerTableBody.innerHTML = '';
    try {
        const response = await authenticatedFetch('/api/providers');
        const data = await response.json();
        if (response.ok) {
            if (data && data.length > 0) {
                data.forEach(provider => {
                    const row = document.createElement('tr');
                    row.classList.add('table-row-hover', 'border-b', 'border-gray-200', 'dark:border-gray-700');
                    row.innerHTML = `
                        <td class="py-3 px-6 text-left">${provider.id}</td>
                        <td class="py-3 px-6 text-left">${provider.name}</td>
                        <td class="py-3 px-6 text-left">${provider.type}</td>
                        <td class="py-3 px-6 text-left">${provider.base_url}</td>
                        <td class="py-3 px-6 text-center">
                            <button class="delete-btn bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-3 rounded focus:outline-none focus:shadow-outline" data-id="${provider.id}" data-type="provider">Delete</button>
                        </td>
                    `;
                    providerTableBody.appendChild(row);
                });
                addDeleteEventListeners();
            } else {
                providerTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">No providers found.</td></tr>`;
            }
        } else {
            const errorData = await response.json();
            providerTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">Error: ${errorData.error || 'Failed to fetch providers'}</td></tr>`;
        }
    } catch (error) {
        console.error('Error fetching providers:', error);
        providerTableBody.innerHTML = `<tr><td colspan="5" class="py-3 px-6 text-center">An error occurred while fetching providers.</td></tr>`;
    }
}


function formatJsonPayload(payload) {
    if (!payload) return 'N/A';
    try {
        const parsed = JSON.parse(payload);
        return `<pre class="json-payload hidden">${JSON.stringify(parsed, null, 2)}</pre><button class="view-json-btn text-blue-500 hover:underline">View Details</button>`;
    } catch (e) {
        return payload.substring(0, 100) + (payload.length > 100 ? '...' : '');
    }
}

// Fetch and display Conversation Logs
export async function fetchConversationLogs() {
    if (!conversationLogTableBody) return; // Ensure element exists
    conversationLogTableBody.innerHTML = '';
    try {
        const response = await authenticatedFetch('/api/conversation_logs');
        const result = await response.json();
        if (response.ok) {
            if (result.logs && result.logs.length > 0) {
                result.logs.forEach(log => {
                    const row = document.createElement('tr');
                    row.classList.add('table-row-hover', 'border-b', 'border-gray-200', 'dark:border-gray-700');
                    row.innerHTML = `
                        <td class="py-3 px-6 text-left">${log.id}</td>
                        <td class="py-3 px-6 text-left">${log.connection_id}</td>
                        <td class="py-3 px-6 text-left">${log.model_id}</td>
                        <td class="py-3 px-6 text-left">${log.prompt_tokens}</td>
                        <td class="py-3 px-6 text-left">${log.completion_tokens}</td>
                        <td class="py-3 px-6 text-left">${formatJsonPayload(log.request_payload)}</td>
                        <td class="py-3 px-6 text-left">${formatJsonPayload(log.response_payload)}</td>
                        <td class="py-3 px-6 text-left">${new Date(log.created_at).toLocaleString()}</td>
                    `;
                    conversationLogTableBody.appendChild(row);
                });
                addJsonViewEventListeners();
            } else {
                conversationLogTableBody.innerHTML = `<tr><td colspan="8" class="py-3 px-6 text-center">No conversation logs found.</td></tr>`;
            }
        } else {
            conversationLogTableBody.innerHTML = `<tr><td colspan="8" class="py-3 px-6 text-center">Error: ${result.error || 'Failed to fetch conversation logs'}</td></tr>`;
        }
    } catch (error) {
        console.error('Error fetching conversation logs:', error);
        conversationLogTableBody.innerHTML = `<tr><td colspan="8" class="py-3 px-6 text-center">An error occurred while fetching conversation logs.</td></tr>`;
    }
}

function addJsonViewEventListeners() {
    document.querySelectorAll('.view-json-btn').forEach(button => {
        button.addEventListener('click', (e) => {
            const payloadPre = e.target.previousElementSibling;
            if (payloadPre && payloadPre.classList.contains('json-payload')) {
                payloadPre.classList.toggle('hidden');
                e.target.textContent = payloadPre.classList.contains('hidden') ? 'View Details' : 'Hide Details';
            }
        });
    });
}

// Add event listeners to delete buttons
function addDeleteEventListeners() {
    document.querySelectorAll('.delete-btn').forEach(button => {
        button.removeEventListener('click', handleDeleteClick); // Prevent duplicate listeners
        button.addEventListener('click', handleDeleteClick);
    });
}

async function handleDeleteClick(e) {
    const id = e.target.dataset.id;
    const type = e.target.dataset.type;
    const result = await handleDelete(id, type);
    if (result.success) {
        // Re-fetch data to update the table
        switch (type) {
            case 'api-key':
                fetchApiKeys();
                break;
            case 'model':
                fetchModels();
                break;
            case 'connection':
                fetchConnections();
                break;
            case 'provider':
                fetchProviders();
                break;
        }
    }
}

export const setupMainUINav = () => {
    // Existing nav button logic for main UI
    navButtons.forEach(button => {
        button.addEventListener('click', () => {
            const sectionId = button.dataset.section + '-section';
            showSection(sectionId);
            activateNavButton(button);

            // Fetch data for the selected section
            switch (button.dataset.section) {
                case 'api-keys':
                    fetchApiKeys();
                    break;
                case 'models':
                    fetchModels();
                    break;
                case 'connections':
                    fetchConnections();
                    break;
                case 'providers':
                    fetchProviders();
                    break;
                case 'conversation-logs':
                    fetchConversationLogs();
                    break;
            }
        });
    });
};

