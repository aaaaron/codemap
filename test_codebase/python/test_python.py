# This is a sample Python file for testing the parser

def greet(name):
    """Greet a person by name."""
    return f"Hello, {name}!"

class Calculator:
    """A simple calculator class."""

    def __init__(self):
        self.result = 0

    def add(self, x, y):
        """Add two numbers."""
        return x + y

    def multiply(self, x, y):
        """Multiply two numbers."""
        return x * y

def main():
    calc = Calculator()
    print(greet("World"))
    print(calc.add(2, 3))