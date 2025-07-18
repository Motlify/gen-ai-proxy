// src/templates/ui/js/main.js
console.log('main.js script started.');

import { setupThemeToggle } from './theme.js';
import { initModal, openModal, closeModal } from './modal.js';
import { createApiKeyFormHtml, createConnectionFormHtml, createProviderFormHtml, createModelFormHtml } from './forms.js';
import { handleCreateApiKey, handleCreateConnection, handleCreateProvider, handleCreateModel, fetchProvidersForSelect, fetchConnectionsForSelect } from './api.js';
import { initAuthElements, showAuthSection, showMainContent, setupAuthForms, setupLogout } from './auth.js';
import { initMainUIElements, showSection, activateNavButton, fetchApiKeys, fetchModels, fetchConnections, fetchProviders, fetchConversationLogs } from './mainUI.js';
import { initDashboard, loadDashboardSection, setupDashboardEventHandlers } from './dashboard.js';


document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const showRegisterLink = document.getElementById('showRegister');
    const showLoginLink = document.getElementById('showLogin');
    const authSection = document.getElementById('auth-section');
    const mainContent = document.getElementById('main-content');
    const logoutButton = document.getElementById('logoutButton');
    const loginMessage = document.getElementById('loginMessage');
    const registerMessage = document.getElementById('registerMessage');
    const themeToggle = document.getElementById('themeToggle');
    const navButtons = document.querySelectorAll('.nav-button');
    const contentSections = document.querySelectorAll('.content-section');
    const apiKeyTableBody = document.getElementById('api-key-table-body');
    const modelTableBody = document.getElementById('model-table-body');
    const connectionTableBody = document.getElementById('connection-table-body');
    const providerTableBody = document.getElementById('provider-table-body');
    const conversationLogTableBody = document.getElementById('conversation-log-table-body');
    const dashboardContent = document.getElementById('dashboardContent');
    const genericModal = document.getElementById('genericModal');
    const closeModalButton = document.getElementById('closeModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalContent = document.getElementById('modalContent');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const authLoginFormContainer = document.getElementById('login-form');
    const authRegisterFormContainer = document.getElementById('register-form');

    // Initialize modules
    setupThemeToggle(themeToggle);
    if (genericModal) {
        initModal(genericModal, closeModalButton, modalTitle, modalContent);
    }
    initAuthElements(authSection, mainContent, loginMessage, registerMessage, authLoginFormContainer, authRegisterFormContainer);
    initMainUIElements(apiKeyTableBody, modelTableBody, connectionTableBody, providerTableBody, conversationLogTableBody, navButtons, contentSections);
    initDashboard(dashboardContent);

    // Check for token on load
    const token = localStorage.getItem('jwt_token');
    if (token) {
        showMainContent(dashboardContent, navButtons, showSection, activateNavButton, fetchApiKeys, loadDashboardSection);
    } else {
        showAuthSection();
    }

    setupAuthForms(loginForm, registerForm, showRegisterLink, showLoginLink, () => showMainContent(dashboardContent, navButtons, showSection, activateNavButton, fetchApiKeys, loadDashboardSection));
    setupLogout(logoutButton);

    if (dashboardContent) { // This means we are on dashboard.html
        loadDashboardSection('api-keys');
    } else { // This means we are on index.html
        navButtons.forEach(button => {
            button.addEventListener('click', () => {
                const sectionId = button.dataset.section + '-section';
                showSection(sectionId);
                activateNavButton(button);
                switch (button.dataset.section) {
                    case 'api-keys': fetchApiKeys(); break;
                    case 'models': fetchModels(); break;
                    case 'connections': fetchConnections(); break;
                    case 'providers': fetchProviders(); break;
                    case 'conversation-logs': fetchConversationLogs(); break;
                }
            });
        });

        // Handle "Add" buttons in the main UI
        document.querySelectorAll('button[data-modal-target]').forEach(button => {
            button.addEventListener('click', async () => {
                const modalTarget = button.dataset.modalTarget;
                let title = '';
                let formHtml = '';
                let submitHandler = null;
                let refreshFunction = null;

                switch (modalTarget) {
                    case 'create-api-key-modal':
                        title = 'Create New API Key';
                        formHtml = createApiKeyFormHtml;
                        submitHandler = handleCreateApiKey;
                        refreshFunction = fetchApiKeys;
                        break;
                    case 'create-model-modal':
                        title = 'Create New Model';
                        formHtml = createModelFormHtml;
                        submitHandler = handleCreateModel;
                        refreshFunction = fetchModels;
                        break;
                    case 'create-connection-modal':
                        title = 'Create New Connection';
                        formHtml = createConnectionFormHtml;
                        submitHandler = handleCreateConnection;
                        refreshFunction = fetchConnections;
                        break;
                    case 'create-provider-modal':
                        title = 'Create New Provider';
                        formHtml = createProviderFormHtml;
                        submitHandler = handleCreateProvider;
                        refreshFunction = fetchProviders;
                        break;
                }

                openModal(title, formHtml);

                if (modalTarget === 'create-provider-modal') {
                    const typeSelect = document.getElementById('provider_type');
                    const baseUrlInput = document.getElementById('provider_base_url');
                    if (typeSelect && baseUrlInput) {
                        typeSelect.addEventListener('change', () => {
                            if (typeSelect.value === 'ollama') {
                                baseUrlInput.placeholder = 'http://localhost:11435';
                            } else if (typeSelect.value === 'openai') {
                                baseUrlInput.placeholder = 'https://api.openai.com/v1';
                            } else {
                                baseUrlInput.placeholder = '';
                            }
                        });
                    }
                }

                if (modalTarget === 'create-connection-modal') {
                    const result = await fetchProvidersForSelect();
                    const selectElement = document.getElementById('connection_provider_id');
                    if (result.success && selectElement) {
                        selectElement.innerHTML = '<option value="">Select a Provider</option>'; // Clear loading text
                        result.data.forEach(provider => {
                            const option = document.createElement('option');
                            option.value = provider.id;
                            option.textContent = `${provider.name} (ID: ${provider.id})`;
                            selectElement.appendChild(option);
                        });
                    } else {
                        selectElement.innerHTML = '<option value="">Failed to load providers</option>';
                    }
                }

                if (modalTarget === 'create-model-modal') {
                    const result = await fetchConnectionsForSelect();
                    const selectElement = document.getElementById('model_connection_id');
                    if (result.success && selectElement) {
                        selectElement.innerHTML = '<option value="">Select a Connection</option>'; // Clear loading text
                        result.data.forEach(connection => {
                            const option = document.createElement('option');
                            option.value = connection.id;
                            option.textContent = `${connection.name} (ID: ${connection.id})`;
                            selectElement.appendChild(option);
                        });
                    } else {
                        selectElement.innerHTML = '<option value="">Failed to load connections</option>';
                    }

                    const modelTypeSelect = document.getElementById('model_type');
                    const thinkingDiv = document.querySelector('div.mb-4:has(#model_thinking)');
                    const toolsUsageDiv = document.querySelector('div.mb-4:has(#model_tools_usage)');

                    const toggleThinkingTools = () => {
                        if (modelTypeSelect.value === 'llm') {
                            thinkingDiv.classList.remove('hidden');
                            toolsUsageDiv.classList.remove('hidden');
                        } else {
                            thinkingDiv.classList.add('hidden');
                            toolsUsageDiv.classList.add('hidden');
                        }
                    };

                    modelTypeSelect.addEventListener('change', toggleThinkingTools);
                    toggleThinkingTools(); // Initial call to set visibility based on default/pre-selected value
                }

                if (submitHandler) {
                    const form = modalContent.querySelector('form');
                    if (form) {
                        form.addEventListener('submit', async (event) => {
                            const result = await submitHandler(event);
                            if (result.success) {
                                closeModal();
                                if (refreshFunction) refreshFunction();
                            }
                        });
                    }
                }
            });
        });
    }
});
