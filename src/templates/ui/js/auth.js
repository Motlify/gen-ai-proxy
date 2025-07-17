// src/templates/ui/js/auth.js

let authSection, mainContent, loginMessage, registerMessage, authLoginFormContainer, authRegisterFormContainer;

export const initAuthElements = (authSect, mainCont, loginMsg, registerMsg, loginFormCont, registerFormCont) => {
    authSection = authSect;
    mainContent = mainCont;
    loginMessage = loginMsg;
    registerMessage = registerMsg;
    authLoginFormContainer = loginFormCont;
    authRegisterFormContainer = registerFormCont;
};

export const showAuthSection = () => {
    if (authSection) authSection.classList.remove('hidden');
    if (mainContent) mainContent.classList.add('hidden');
};

export const showMainContent = (dashboardContent, navButtons, showSection, activateNavButton, fetchApiKeys, loadDashboardSectionCallback) => {
    if (authSection) authSection.classList.add('hidden');
    if (mainContent) mainContent.classList.remove('hidden');

    // Prioritize dashboard content if dashboardContent element exists
    if (dashboardContent) {
        loadDashboardSectionCallback('api-keys'); 
    } else if (navButtons.length > 0) {
        // Default to API Keys section on login for main UI (index.html)
        showSection('api-keys-section');
        activateNavButton(document.querySelector('[data-section="api-keys"]'));
        fetchApiKeys(); // Fetch API keys for the default view in main UI
    }
};

export const setupAuthForms = (loginForm, registerForm, showRegisterLink, showLoginLink, showMainContentCallback) => {
    // Toggle between login and register forms
    if (showRegisterLink) {
        showRegisterLink.addEventListener('click', (e) => {
            e.preventDefault();
            authLoginFormContainer.classList.add('hidden');
            authRegisterFormContainer.classList.remove('hidden');
            loginMessage.textContent = ''; // Clear messages
            registerMessage.textContent = '';
        });
    }

    if (showLoginLink) {
        showLoginLink.addEventListener('click', (e) => {
            e.preventDefault();
            authRegisterFormContainer.classList.add('hidden');
            authLoginFormContainer.classList.remove('hidden');
            loginMessage.textContent = ''; // Clear messages
            registerMessage.textContent = '';
        });
    }

    // Handle Login
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = e.target.loginUsername.value;
            const password = e.target.loginPassword.value;

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password }),
                });

                const data = await response.json();

                if (response.ok) {
                    localStorage.setItem('jwt_token', data.access_token);
                    loginMessage.textContent = '';
                    showMainContentCallback();
                } else {
                    loginMessage.textContent = data.error || 'Login failed';
                    console.error('Login failed:', data.error || response.statusText);
                }
            } catch (error) {
                console.error('Error during login fetch:', error);
                loginMessage.textContent = 'An error occurred. Please try again.';
            }
        });
    }

    // Handle Registration
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = e.target.registerUsername.value;
            const password = e.target.registerPassword.value;

            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password }),
                });

                const data = await response.json();

                if (response.ok) {
                    registerMessage.textContent = 'Registration successful! You can now log in.';
                    registerMessage.style.color = 'green';
                    authLoginFormContainer.classList.remove('hidden');
                    authRegisterFormContainer.classList.add('hidden');
                } else {
                    registerMessage.textContent = data.error || 'Registration failed';
                    registerMessage.style.color = 'red';
                }
            } catch (error) {
                console.error('Error during registration:', error);
                registerMessage.textContent = 'An error occurred. Please try again.';
                registerMessage.style.color = 'red';
            }
        });
    }
};

export const setupLogout = (logoutButton) => {
    // Handle Logout
    if (logoutButton) {
        logoutButton.addEventListener('click', (e) => {
            e.preventDefault(); // Prevent default anchor behavior
            localStorage.removeItem('jwt_token');
            // Redirect to the root or login page, depending on context
            if (window.location.pathname === '/dashboard.html') {
                window.location.href = '/'; // Redirect from dashboard to main login
            } else {
                showAuthSection(); // Stay on main page, show auth
            }
        });
    }
};
