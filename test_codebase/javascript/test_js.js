// Test JavaScript file for parser testing

/**
 * Greets a person by name
 * @param {string} name - The person's name
 * @returns {string} The greeting message
 */
function greet(name) {
    return `Hello, ${name}!`;
}

/**
 * Calculator class for basic arithmetic operations
 */
class Calculator {
    constructor() {
        this.result = 0;
    }

    /**
     * Adds two numbers
     * @param {number} x - First number
     * @param {number} y - Second number
     * @returns {number} The sum
     */
    add(x, y) {
        return x + y;
    }

    /**
     * Multiplies two numbers
     * @param {number} x - First number
     * @param {number} y - Second number
     * @returns {number} The product
     */
    multiply(x, y) {
        return x * y;
    }
}

// Arrow function example
const sayHello = (name) => `Hello, ${name}!`;

export { greet, Calculator, sayHello };