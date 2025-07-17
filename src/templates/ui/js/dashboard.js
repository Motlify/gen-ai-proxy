// src/templates/ui/js/dashboard.js

import { authenticatedFetch, handleCreateApiKey, handleCreateConnection, handleCreateProvider, handleCreateModel, handleDelete } from './api.js';
import { openModal, closeModal } from './modal.js';
import { createApiKeyFormHtml, createConnectionFormHtml, createProviderFormHtml, createModelFormHtml } from './forms.js';

let dashboardContent;

export const initDashboard = (dashboardContentElement) => {
    dashboardContent = dashboardContentElement;
};

// Function to load content into the dashboard
export const loadDashboardSection = async (section) => {
    if (dashboardContent) {
        dashboardContent.innerHTML = ''; // Clear current content

        let content = '';
        switch (section) {
            case 'api-keys':
                content = `
                    <h2 class="text-2xl font-bold mb-4">API Keys</h2>
                    <div id="apiKeysList">
                        <h3 class="text-xl font-semibold mb-2">Existing API Keys</h3>
                        <p>Loading API keys...</p>
                    </div>
                `;
                break;
            case 'connections':
                content = `
                    <h2 class="text-2xl font-bold mb-4">Connections</h2>
                    <div id="connectionsList">
                        <h3 class="text-xl font-semibold mb-2">Existing Connections</h3>
                        <p>Loading connections...</p>
                    </div>
                `;
                break;
            case 'providers':
                content = `
                    <h2 class="text-2xl font-bold mb-4">Providers</h2>
                    <div id="providersList">
                        <h3 class="text-xl font-semibold mb-2">Existing Providers</h3>
                        <p>Loading providers...</p>
                    </div>
                `;
                break;
            case 'models':
                content = `
                    <h2 class="text-2xl font-bold mb-4">Models</h2>
                    <div id="modelsList">
                        <h3 class="text-xl font-semibold mb-2">Existing Models</h3>
                        <p>Loading models...</p>
                    </div>
                `;
                break;
            case 'conversation-logs':
                content = `
                    <h2 class="text-2xl font-bold mb-4">Conversation Logs</h2>
                    <div id="conversationLogsList">
                        <h3 class="text-xl font-semibold mb-2">Existing Conversation Logs</h3>
                        <p>Loading conversation logs...</p>
                    </div>
                `;
                break;
            default:
                content = `<h1 class="text-3xl font-bold text-center mb-4">Welcome to the GenAI Proxy Dashboard!</h1><p class="text-center text-gray-400">Use the navigation above to manage your resources.</p>`;
        }
        dashboardContent.innerHTML = content;

        // Fetch data for the selected section
        if (section === 'api-keys') {
            await fetchApiKeysDashboard();
        } else if (section === 'connections') {
            await fetchConnectionsDashboard();
        } else if (section === 'providers') {
            await fetchProvidersDashboard();
        } else if (section === 'models') {
            await fetchModelsDashboard();
        } else if (section === 'conversation-logs') {
            await fetchConversationLogsDashboard();
        }
        // Re-attach event handlers after content is loaded
        setupDashboardEventHandlers(document.querySelectorAll('nav a[data-section]'), document.querySelectorAll('button[data-modal-target]'));
    }
};

// Fetch and display API Keys (dashboard-specific)
export const fetchApiKeysDashboard = async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
        const apiKeysListDiv = document.getElementById('apiKeysList');
        if (apiKeysListDiv) apiKeysListDiv.innerHTML = '<p class="text-red-400">Please log in to view API keys.</p>';
        return;
    }

    try {
        const response = await authenticatedFetch('/api/api-keys');
        const data = await response.json();
        const apiKeysListDiv = document.getElementById('apiKeysList');

        if (response.ok) {
            if (apiKeysListDiv) {
                if (data.api_keys && data.api_keys.length > 0) {
                    const ul = document.createElement('ul');
                    ul.className = 'list-disc pl-5';
                    data.api_keys.forEach(key => {
                        const li = document.createElement('li');
                        li.className = 'mb-2';
                        li.innerHTML = `<strong>Name:</strong> ${key.name} (ID: ${key.id}) - Created: ${new Date(key.created_at).toLocaleString()} ${key.last_used_at ? `(Last Used: ${new Date(key.last_used_at).toLocaleString()})` : ''}`;
                        ul.appendChild(li);
                    });
                    apiKeysListDiv.innerHTML = '';
                    apiKeysListDiv.appendChild(ul);
                } else {
                    apiKeysListDiv.innerHTML = '<p>No API keys found.</p>';
                }
            }
        } else {
            if (apiKeysListDiv) apiKeysListDiv.innerHTML = `<p class="text-red-400">Error loading API keys: ${data.error}</p>`;
        }
    } catch (error) {
        console.error('Error fetching API keys:', error);
        const apiKeysListDiv = document.getElementById('apiKeysList');
        if (apiKeysListDiv) apiKeysListDiv.innerHTML = '<p class="text-red-400">An error occurred while fetching API keys.</p>';
    }
};

// Fetch and display Connections (dashboard-specific)
export const fetchConnectionsDashboard = async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
        const connectionsListDiv = document.getElementById('connectionsList');
        if (connectionsListDiv) connectionsListDiv.innerHTML = '<p class="text-red-400">Please log in to view connections.</p>';
        return;
    }

    try {
        const response = await authenticatedFetch('/api/connections');
        const data = await response.json();
        const connectionsListDiv = document.getElementById('connectionsList');

        if (response.ok) {
            if (connectionsListDiv) {
                if (data.connections && data.connections.length > 0) {
                    const ul = document.createElement('ul');
                    ul.className = 'list-disc pl-5';
                    data.connections.forEach(conn => {
                        const li = document.createElement('li');
                        li.className = 'mb-2';
                        li.innerHTML = `<strong>Name:</strong> ${conn.name} (ID: ${conn.id}) - Provider: ${conn.provider} - Created: ${new Date(conn.created_at).toLocaleString()}`;
                        ul.appendChild(li);
                    });
                    connectionsListDiv.innerHTML = '';
                    connectionsListDiv.appendChild(ul);
                } else {
                    connectionsListDiv.innerHTML = '<p>No connections found.</p>';
                }
            }
        } else {
            if (connectionsListDiv) connectionsListDiv.innerHTML = `<p class="text-red-400">Error loading connections: ${data.error}</p>`;
        }
    } catch (error) {
        console.error('Error fetching connections:', error);
        const connectionsListDiv = document.getElementById('connectionsList');
        if (connectionsListDiv) connectionsListDiv.innerHTML = '<p class="text-red-400">An error occurred while fetching connections.</p>';
    }
};

// Fetch and display Providers (dashboard-specific)
export const fetchProvidersDashboard = async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
        const providersListDiv = document.getElementById('providersList');
        if (providersListDiv) providersListDiv.innerHTML = '<p class="text-red-400">Please log in to view providers.</p>';
        return;
    }

    try {
        const response = await authenticatedFetch('/api/providers');
        const data = await response.json();
        const providersListDiv = document.getElementById('providersList');

        if (response.ok) {
            if (providersListDiv) {
                if (data && data.length > 0) {
                    const ul = document.createElement('ul');
                    ul.className = 'list-disc pl-5';
                    data.forEach(provider => {
                        const li = document.createElement('li');
                        li.className = 'mb-2';
                        li.innerHTML = `<strong>Name:</strong> ${provider.name} (ID: ${provider.id}) - Type: ${provider.type} - Base URL: ${provider.base_url || 'N/A'}`;
                        ul.appendChild(li);
                    });
                    providersListDiv.innerHTML = '';
                    providersListDiv.appendChild(ul);
                } else {
                    providersListDiv.innerHTML = '<p>No providers found.</p>';
                }
            }
        } else {
            if (providersListDiv) providersListDiv.innerHTML = `<p class="text-red-400">Error loading providers: ${data.error}</p>`;
        }
    } catch (error) {
        console.error('Error fetching providers:', error);
        const providersListDiv = document.getElementById('providersList');
        if (providersListDiv) providersListDiv.innerHTML = '<p class="text-red-400">An error occurred while fetching providers.</p>';
    }
};

// Fetch and display Models (dashboard-specific)
export const fetchModelsDashboard = async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
        const modelsListDiv = document.getElementById('modelsList');
        if (modelsListDiv) modelsListDiv.innerHTML = '<p class="text-red-400">Please log in to view models.</p>';
        return;
    }

    try {
        const response = await authenticatedFetch('/api/models');
        const data = await response.json();
        const modelsListDiv = document.getElementById('modelsList');

        if (response.ok) {
            if (modelsListDiv) {
                if (data && data.length > 0) {
                    const ul = document.createElement('ul');
                    ul.className = 'list-disc pl-5';
                    data.forEach(model => {
                        const li = document.createElement('li');
                        li.className = 'mb-2';
                        li.innerHTML = `<strong>Proxy Model ID:</strong> ${model.proxy_model_id} (ID: ${model.id}) - Provider Model ID: ${model.provider_model_id} - Connection ID: ${model.connection_id} - Price Input: ${model.price_input} - Price Output: ${model.price_output} - Thinking: ${model.thinking} - Tools Usage: ${model.tools_usage}`;
                        ul.appendChild(li);
                    });
                    modelsListDiv.innerHTML = '';
                    modelsListDiv.appendChild(ul);
                } else {
                    modelsListDiv.innerHTML = '<p>No models found.</p>';
                }
            }
        } else {
            if (modelsListDiv) modelsListDiv.innerHTML = `<p class="text-red-400">Error loading models: ${data.error}</p>`;
        }
    } catch (error) {
        console.error('Error fetching models:', error);
        const modelsListDiv = document.getElementById('modelsList');
        if (modelsListDiv) modelsListDiv.innerHTML = '<p class="text-red-400">An error occurred while fetching models.</p>';
    }
};

function formatJsonPayload(payload) {
    if (!payload) return 'N/A';
    try {
        const parsed = JSON.parse(payload);
        return `<pre class="json-payload hidden">${JSON.stringify(parsed, null, 2)}</pre><button class="view-json-btn text-blue-500 hover:underline">View Details</button>`;
    } catch (e) {
        return payload.substring(0, 100) + (payload.length > 100 ? '...' : '');
    }
}

// Fetch and display Conversation Logs (dashboard-specific)
export const fetchConversationLogsDashboard = async () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) {
        const conversationLogsListDiv = document.getElementById('conversationLogsList');
        if (conversationLogsListDiv) conversationLogsListDiv.innerHTML = '<p class="text-red-400">Please log in to view conversation logs.</p>';
        return;
    }

    try {
        const response = await authenticatedFetch('/api/conversation_logs');
        const data = await response.json();
        const conversationLogsListDiv = document.getElementById('conversationLogsList');

        if (response.ok) {
            if (conversationLogsListDiv) {
                if (data.logs && data.logs.length > 0) {
                    const ul = document.createElement('ul');
                    ul.className = 'list-disc pl-5';
                    data.logs.forEach(log => {
                        const li = document.createElement('li');
                        li.className = 'mb-2';
                        li.innerHTML = `<strong>ID:</strong> ${log.id} - Connection ID: ${log.connection_id} - Model ID: ${log.model_id} - Prompt Tokens: ${log.prompt_tokens} - Completion Tokens: ${log.completion_tokens} - Request Payload: ${formatJsonPayload(log.request_payload)} - Response Payload: ${formatJsonPayload(log.response_payload)} - Created At: ${new Date(log.created_at).toLocaleString()}`;
                        ul.appendChild(li);
                    });
                    conversationLogsListDiv.innerHTML = '';
                    conversationLogsListDiv.appendChild(ul);
                    addJsonViewEventListenersDashboard();
                } else {
                    conversationLogsListDiv.innerHTML = '<p>No conversation logs found.</p>';
                }
            }
        } else {
            if (conversationLogsListDiv) conversationLogsListDiv.innerHTML = `<p class="text-red-400">Error loading conversation logs: ${data.error}</p>`;
        }
    } catch (error) {
        console.error('Error fetching conversation logs:', error);
        const conversationLogsListDiv = document.getElementById('conversationLogsList');
        if (conversationLogsListDiv) conversationLogsListDiv.innerHTML = '<p class="text-red-400">An error occurred while fetching conversation logs.</p>';
    }
};

function addJsonViewEventListenersDashboard() {
    document.querySelectorAll('#conversationLogsList .view-json-btn').forEach(button => {
        button.addEventListener('click', (e) => {
            const payloadPre = e.target.previousElementSibling;
            if (payloadPre && payloadPre.classList.contains('json-payload')) {
                payloadPre.classList.toggle('hidden');
                e.target.textContent = payloadPre.classList.contains('hidden') ? 'View Details' : 'Hide Details';
            }
        });
    });
}

export const setupDashboardEventHandlers = (dashboardNavLinks, addButtons) => {
    console.log('setupDashboardEventHandlers called.');
    console.log('Number of nav links:', dashboardNavLinks.length);
    console.log('Number of add buttons:', addButtons.length);
    if (dashboardContent) { // This means we are on dashboard.html
        dashboardNavLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const section = e.target.dataset.section;
                console.log('Nav link clicked:', section);
                loadDashboardSection(section);
            });
        });

        addButtons.forEach(button => {
            console.log('Attaching click listener to button:', button.dataset.modalTarget);
            button.addEventListener('click', async (e) => {
                console.log('Button clicked:', e.target.dataset.modalTarget);
                const modalTarget = e.target.dataset.modalTarget;
                let title = '';
                let formHtml = '';
                let submitHandler = null;
                let refreshFunction = null;

                switch (modalTarget) {
                    case 'create-api-key-modal':
                        title = 'Create New API Key';
                        formHtml = createApiKeyFormHtml;
                        submitHandler = handleCreateApiKey;
                        refreshFunction = fetchApiKeysDashboard;
                        break;
                    case 'create-connection-modal':
                        title = 'Create New Connection';
                        formHtml = createConnectionFormHtml;
                        submitHandler = handleCreateConnection;
                        refreshFunction = fetchConnectionsDashboard;
                        break;
                    case 'create-provider-modal':
                        title = 'Create New Provider';
                        formHtml = createProviderFormHtml;
                        submitHandler = handleCreateProvider;
                        refreshFunction = fetchProvidersDashboard;
                        break;
                    case 'create-model-modal':
                        console.log('Case: create-model-modal');
                        title = 'Create New Model';
                        formHtml = createModelFormHtml;
                        submitHandler = handleCreateModel;
                        refreshFunction = fetchModelsDashboard;
                        break;
                    default:
                        console.log('Unknown modal target:', modalTarget);
                        return;
                }
                console.log('Opening modal with title:', title);
                openModal(title, formHtml);

                if (submitHandler) {
                    const form = modalContent.querySelector('form');
                    if (form) {
                        console.log('Attaching submit handler to form:', form.id);
                        form.addEventListener('submit', async (event) => {
                            console.log('Form submitted:', form.id);
                            const result = await submitHandler(event);
                            if (result.success) {
                                closeModal();
                                if (refreshFunction) refreshFunction();
                            }
                        });
                    } else {
                        console.error('Form element not found for modal target:', modalTarget);
                    }
                }
            });
        });

        // Load default section on dashboard page load (e.g., API Keys)
        console.log('Loading default dashboard section: api-keys');
        loadDashboardSection('api-keys');
    }
};