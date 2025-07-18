// src/templates/ui/js/forms.js

export const createApiKeyFormHtml = `
    <form id="createApiKeyForm" class="bg-gray-700 p-4 rounded-lg shadow-md">
        <div class="mb-4">
            <label for="api_key_name" class="block text-sm font-medium text-gray-300">Name</label>
            <input type="text" id="api_key_name" name="name" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Create API Key
        </button>
    </form>
`;

export const createConnectionFormHtml = `
    <form id="createConnectionForm" class="bg-gray-700 p-4 rounded-lg shadow-md">
        <div class="mb-4">
            <label for="connection_name" class="block text-sm font-medium text-gray-300">Name</label>
            <input type="text" id="connection_name" name="name" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4">
            <label for="connection_provider_id" class="block text-sm font-medium text-gray-300">Provider</label>
            <select id="connection_provider_id" name="provider_id" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option value="">Loading providers...</option>
            </select>
        </div>
        <div class="mb-4">
            <label for="connection_api_key" class="block text-sm font-medium text-gray-300">API Key</label>
            <input type="text" id="connection_api_key" name="api_key" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
            <p class="mt-2 text-xs text-gray-400">This is the API key for the selected provider.</p>
        </div>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Create Connection
        </button>
    </form>
`;

export const createProviderFormHtml = `
    <form id="createProviderForm" class="bg-gray-700 p-4 rounded-lg shadow-md">
        <div class="mb-4">
            <label for="provider_name" class="block text-sm font-medium text-gray-300">Name</label>
            <input type="text" id="provider_name" name="name" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4">
            <label for="provider_type" class="block text-sm font-medium text-gray-300">Type</label>
            <select id="provider_type" name="type" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option value="">Select a type</option>
                <option value="ollama">Ollama</option>
                <option value="openai">OpenAI (compatible)</option>
            </select>
        </div>
        <div class="mb-4">
            <label for="provider_base_url" class="block text-sm font-medium text-gray-300">Base URL</label>
            <input type="url" id="provider_base_url" name="base_url"
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Create Provider
        </button>
    </form>
`;

export const createModelFormHtml = `
    <form id="createModelForm" class="bg-gray-700 p-4 rounded-lg shadow-md">
        <div class="mb-4">
            <label for="model_connection_id" class="block text-sm font-medium text-gray-300">Connection</label>
            <select id="model_connection_id" name="connection_id" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option value="">Loading connections...</option>
            </select>
        </div>
        <div class="mb-4">
            <label for="model_proxy_model_id" class="block text-sm font-medium text-gray-300">Proxy Model ID</label>
            <input type="text" id="model_proxy_model_id" name="proxy_model_id" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4">
            <label for="model_provider_model_id" class="block text-sm font-medium text-gray-300">Provider Model ID</label>
            <input type="text" id="model_provider_model_id" name="provider_model_id" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4">
            <label for="model_type" class="block text-sm font-medium text-gray-300">Model Type</label>
            <select id="model_type" name="type" required
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option value="">Select a type</option>
                <option value="llm">LLM</option>
                <option value="embedding">Embedding</option>
            </select>
        </div>
        <div class="mb-4">
            <label for="model_price_input" class="block text-sm font-medium text-gray-300">Price Input</label>
            <input type="number" step="0.000001" id="model_price_input" name="price_input"
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4">
            <label for="model_price_output" class="block text-sm font-medium text-gray-300">Price Output</label>
            <input type="number" step="0.000001" id="model_price_output" name="price_output"
                   class="mt-1 block w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
        </div>
        <div class="mb-4 flex items-center">
            <input type="checkbox" id="model_thinking" name="thinking"
                   class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-600 rounded">
            <label for="model_thinking" class="ml-2 block text-sm text-gray-300">Thinking</label>
        </div>
        <div class="mb-4 flex items-center">
            <input type="checkbox" id="model_tools_usage" name="tools_usage"
                   class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-600 rounded">
            <label for="model_tools_usage" class="ml-2 block text-sm text-gray-300">Tools Usage</label>
        </div>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Create Model
        </button>
    </form>
`;