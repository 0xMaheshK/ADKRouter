export class ADKRouterError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'ADKRouterError';
    }
}

export class ConfigurationError extends ADKRouterError {
    constructor(message: string) {
        super(message);
        this.name = 'ConfigurationError';
    }
}

export class ProviderError extends ADKRouterError {
    constructor(message: string) {
        super(message);
        this.name = 'ProviderError';
    }
}

export class ModelError extends ADKRouterError {
    constructor(message: string) {
        super(message);
        this.name = 'ModelError';
    }
}