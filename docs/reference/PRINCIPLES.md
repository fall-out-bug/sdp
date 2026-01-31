# Software Engineering Principles

Core principles that guide all development under SDP. These are **non-negotiable**.

---

## Quick Reference

| Principle | One-liner | Violation Example |
|-----------|-----------|-------------------|
| **SRP** | One class = one reason to change | `UserService` does auth, email, and logging |
| **OCP** | Extend, don't modify | Adding `if type == "new"` everywhere |
| **LSP** | Subtypes must be substitutable | `Square` breaks `Rectangle.setWidth()` |
| **ISP** | Small, focused interfaces | `IWorker` with `work()` + `eat()` |
| **DIP** | Depend on abstractions | `OrderService` imports `MySQLDatabase` |
| **DRY** | Don't repeat yourself | Same validation logic in 5 places |
| **KISS** | Keep it simple | Regex for simple string check |
| **YAGNI** | Build only what's needed | Adding "future" config options |
| **TDD** | Tests first, then code | Writing tests after implementation |

---

## SOLID Principles

### S — Single Responsibility Principle (SRP)

**A class should have only one reason to change.**

```python
# BAD: Multiple responsibilities
class UserService:
    def authenticate(self, email: str, password: str) -> User: ...
    def send_welcome_email(self, user: User) -> None: ...
    def generate_report(self, user: User) -> str: ...
    def log_activity(self, user: User, action: str) -> None: ...

# GOOD: Separated concerns
class AuthService:
    def authenticate(self, email: str, password: str) -> User: ...

class EmailService:
    def send_welcome(self, user: User) -> None: ...

class ReportGenerator:
    def generate_user_report(self, user: User) -> str: ...

class ActivityLogger:
    def log(self, user: User, action: str) -> None: ...
```

**AI Prompt**: "Does this class have more than one reason to change? If yes, split it."

---

### O — Open/Closed Principle (OCP)

**Open for extension, closed for modification.**

```python
# BAD: Modifying existing code for new types
class PaymentProcessor:
    def process(self, payment_type: str, amount: float) -> None:
        if payment_type == "credit_card":
            self._process_credit_card(amount)
        elif payment_type == "paypal":
            self._process_paypal(amount)
        elif payment_type == "crypto":  # Adding new type = modifying class
            self._process_crypto(amount)

# GOOD: Extending via new classes
class PaymentProcessor(Protocol):
    def process(self, amount: float) -> None: ...

class CreditCardProcessor:
    def process(self, amount: float) -> None: ...

class PayPalProcessor:
    def process(self, amount: float) -> None: ...

class CryptoProcessor:  # New type = new class, no modification
    def process(self, amount: float) -> None: ...
```

**AI Prompt**: "Can I add this feature without modifying existing code?"

---

### L — Liskov Substitution Principle (LSP)

**Subtypes must be substitutable for their base types.**

```python
# BAD: Subtype breaks parent contract
class Rectangle:
    def set_width(self, width: int) -> None:
        self.width = width

    def set_height(self, height: int) -> None:
        self.height = height

class Square(Rectangle):  # Violates LSP!
    def set_width(self, width: int) -> None:
        self.width = width
        self.height = width  # Unexpected side effect

    def set_height(self, height: int) -> None:
        self.width = height  # Unexpected side effect
        self.height = height

# GOOD: Proper abstraction
class Shape(Protocol):
    def area(self) -> float: ...

class Rectangle:
    def __init__(self, width: float, height: float):
        self.width = width
        self.height = height

    def area(self) -> float:
        return self.width * self.height

class Square:
    def __init__(self, side: float):
        self.side = side

    def area(self) -> float:
        return self.side ** 2
```

**AI Prompt**: "Can I replace the parent class with this subclass without breaking anything?"

---

### I — Interface Segregation Principle (ISP)

**Clients should not depend on interfaces they don't use.**

```python
# BAD: Fat interface
class IWorker(Protocol):
    def work(self) -> None: ...
    def eat(self) -> None: ...
    def sleep(self) -> None: ...

class Robot:  # Robots don't eat or sleep!
    def work(self) -> None: ...
    def eat(self) -> None:
        raise NotImplementedError  # Violation!
    def sleep(self) -> None:
        raise NotImplementedError  # Violation!

# GOOD: Segregated interfaces
class IWorkable(Protocol):
    def work(self) -> None: ...

class IFeedable(Protocol):
    def eat(self) -> None: ...

class Human:
    def work(self) -> None: ...
    def eat(self) -> None: ...

class Robot:
    def work(self) -> None: ...  # Only implements what it needs
```

**AI Prompt**: "Does this interface have methods that some implementations won't use?"

---

### D — Dependency Inversion Principle (DIP)

**Depend on abstractions, not concretions.**

```python
# BAD: High-level depends on low-level
class OrderService:
    def __init__(self):
        self.db = MySQLDatabase()  # Direct dependency on concrete class
        self.emailer = SmtpEmailer()  # Direct dependency

    def create_order(self, order: Order) -> None:
        self.db.save(order)
        self.emailer.send_confirmation(order)

# GOOD: Depend on abstractions
class OrderRepository(Protocol):
    def save(self, order: Order) -> None: ...

class EmailSender(Protocol):
    def send_confirmation(self, order: Order) -> None: ...

class OrderService:
    def __init__(
        self,
        repository: OrderRepository,  # Abstract
        emailer: EmailSender  # Abstract
    ):
        self.repository = repository
        self.emailer = emailer

    def create_order(self, order: Order) -> None:
        self.repository.save(order)
        self.emailer.send_confirmation(order)
```

**AI Prompt**: "Am I importing concrete implementations or abstract interfaces?"

---

## DRY — Don't Repeat Yourself

**Every piece of knowledge must have a single, unambiguous representation.**

```python
# BAD: Repeated validation
def create_user(email: str) -> User:
    if "@" not in email or "." not in email:
        raise InvalidEmailError(email)
    ...

def update_email(user: User, email: str) -> None:
    if "@" not in email or "." not in email:  # Duplicated!
        raise InvalidEmailError(email)
    ...

def send_invite(email: str) -> None:
    if "@" not in email or "." not in email:  # Duplicated!
        raise InvalidEmailError(email)
    ...

# GOOD: Single source of truth
class Email:
    def __init__(self, value: str):
        if not self._is_valid(value):
            raise InvalidEmailError(value)
        self.value = value

    @staticmethod
    def _is_valid(value: str) -> bool:
        return "@" in value and "." in value

def create_user(email: Email) -> User: ...
def update_email(user: User, email: Email) -> None: ...
def send_invite(email: Email) -> None: ...
```

**AI Prompt**: "Is this logic duplicated elsewhere? Extract to a single place."

---

## KISS — Keep It Simple, Stupid

**The simplest solution is usually the best.**

```python
# BAD: Over-engineered
def is_palindrome(s: str) -> bool:
    import re
    cleaned = re.sub(r'[^a-zA-Z0-9]', '', s).lower()
    stack = []
    for char in cleaned:
        stack.append(char)
    reversed_str = ''
    while stack:
        reversed_str += stack.pop()
    return cleaned == reversed_str

# GOOD: Simple and clear
def is_palindrome(s: str) -> bool:
    cleaned = ''.join(c.lower() for c in s if c.isalnum())
    return cleaned == cleaned[::-1]
```

**AI Prompt**: "Is there a simpler way to achieve this?"

---

## YAGNI — You Ain't Gonna Need It

**Don't build features until they're actually needed.**

```python
# BAD: Building for hypothetical future
class Config:
    def __init__(
        self,
        database_url: str,
        cache_url: str | None = None,  # "Might need cache later"
        message_queue_url: str | None = None,  # "Might need queue later"
        feature_flags: dict[str, bool] | None = None,  # "For A/B tests"
        plugin_directory: str | None = None,  # "For extensibility"
    ):
        ...

# GOOD: Build only what's needed NOW
class Config:
    def __init__(self, database_url: str):
        self.database_url = database_url

# Add cache_url when you actually need caching
# Add message_queue_url when you actually need a queue
```

**AI Prompt**: "Do I need this right now, or am I building for a hypothetical future?"

---

## TDD — Test-Driven Development

**Write tests before code. Red → Green → Refactor.**

### The Cycle

```
1. RED:    Write a failing test
2. GREEN:  Write minimal code to pass
3. REFACTOR: Improve code, tests still pass
```

### Example

```python
# Step 1: RED - Write failing test
def test_user_can_be_created():
    user = User(email="test@example.com", name="Test")
    assert user.email == "test@example.com"
    assert user.name == "Test"
# Result: NameError: name 'User' is not defined

# Step 2: GREEN - Minimal implementation
@dataclass
class User:
    email: str
    name: str
# Result: Test passes

# Step 3: REFACTOR - Add validation
@dataclass
class User:
    email: str
    name: str

    def __post_init__(self) -> None:
        if "@" not in self.email:
            raise ValueError("Invalid email")
# Result: Tests still pass, code improved
```

### Benefits

1. **Design**: Tests force you to think about API first
2. **Documentation**: Tests document expected behavior
3. **Confidence**: Refactor without fear
4. **Coverage**: 100% coverage by design

### Anti-patterns

```python
# BAD: Tests after implementation
class Calculator:
    def add(self, a, b):
        return a + b

# ... later ...
def test_add():  # Afterthought
    assert Calculator().add(1, 2) == 3

# GOOD: Test first
def test_calculator_adds_two_numbers():
    calc = Calculator()
    assert calc.add(2, 3) == 5
# Now implement Calculator to make test pass
```

**AI Prompt**: "Write the test first, then implement the minimum code to pass it."

---

## Clean Code

### Meaningful Names

```python
# BAD
def calc(d: list[dict]) -> float:
    t = 0
    for i in d:
        t += i['a'] * i['p']
    return t

# GOOD
def calculate_order_total(line_items: list[LineItem]) -> float:
    total = 0.0
    for item in line_items:
        total += item.quantity * item.price
    return total
```

### Small Functions

```python
# BAD: 50-line function doing multiple things
def process_order(order: Order) -> None:
    # validate (10 lines)
    # calculate totals (10 lines)
    # apply discounts (10 lines)
    # update inventory (10 lines)
    # send notification (10 lines)

# GOOD: Composed of focused functions
def process_order(order: Order) -> None:
    validate_order(order)
    totals = calculate_totals(order)
    discounted = apply_discounts(totals)
    update_inventory(order)
    notify_customer(order)
```

### No Comments for Bad Code

```python
# BAD: Comment explaining confusing code
# Check if user is admin by looking at role field which is 1 for admin
if user.role == 1:
    ...

# GOOD: Self-documenting code
if user.is_admin():
    ...
```

---

## Clean Architecture

Dependencies point **inward**. Inner layers know nothing about outer layers.

```
┌─────────────────────────────────────────────────────┐
│                  Presentation                        │
│              (Controllers, Views, API)               │
├─────────────────────────────────────────────────────┤
│                  Infrastructure                      │
│          (Database, External APIs, Frameworks)       │
├─────────────────────────────────────────────────────┤
│                   Application                        │
│              (Use Cases, Services)                   │
├─────────────────────────────────────────────────────┤
│                     Domain                           │
│           (Entities, Business Rules)                 │
└─────────────────────────────────────────────────────┘

        ↑ Dependencies point INWARD (toward Domain)
```

**Layer Rules:**

| Layer | Can Import From | Cannot Import From |
|-------|-----------------|-------------------|
| Domain | Nothing | Application, Infrastructure, Presentation |
| Application | Domain | Infrastructure, Presentation |
| Infrastructure | Domain, Application | Presentation |
| Presentation | Application | - |

**Detailed guide**: [docs/concepts/clean-architecture/README.md](concepts/clean-architecture/README.md)

---

## Verification Checklist

Before completing any workstream:

- [ ] **SRP**: Each class has one responsibility
- [ ] **OCP**: New features don't modify existing code
- [ ] **LSP**: Subtypes are substitutable
- [ ] **ISP**: Interfaces are minimal and focused
- [ ] **DIP**: Dependencies are abstract, not concrete
- [ ] **DRY**: No duplicated logic
- [ ] **KISS**: Simplest solution chosen
- [ ] **YAGNI**: No speculative features
- [ ] **TDD**: Tests written before implementation
- [ ] **Clean Code**: Readable, small functions, good names
- [ ] **Clean Architecture**: No layer violations

---

## References

- [Clean Architecture](concepts/clean-architecture/README.md) — Detailed layer guide
- [Code Patterns](../CODE_PATTERNS.md) — Implementation patterns
- [Protocol](../PROTOCOL.md) — SDP guardrails and gates

---

**Version:** 1.0
**Last Updated:** 2026-01-12
