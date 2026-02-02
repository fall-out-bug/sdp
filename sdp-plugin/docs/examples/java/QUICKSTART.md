# Java Quick Start

SDP workflow for Java projects with Maven/Gradle, JaCoCo, and checkstyle.

## Prerequisites

```bash
# Java 17+
java -version

# Maven 3.8+ or Gradle 8+
mvn --version
# OR
gradle --version
```

## Project Structure

```
my-project/
├── src/
│   ├── main/
│   │   └── java/
│   │       └── com/
│   │           └── example/
│   │               ├── service/
│   │               └── model/
│   └── test/
│       └── java/
│           └── com/
│               └── example/
│                   └── service/
├── pom.xml  # or build.gradle
└── .claude/
    └── skills/
```

## Workflow

### 1. Initialize Project

```bash
# Create feature
@feature "Add user authentication"
```

Claude will interview you about:
- Mission and users
- Technical approach (Maven vs Gradle)
- Success criteria
- Tradeoffs

### 2. Plan Workstreams

```bash
@design feature-auth
```

Claude will:
- Explore Java project structure
- Design workstream decomposition
- Define dependencies
- Request approval

### 3. Execute Workstream

```bash
@build 00-001-01
```

SDP will:
1. Detect Java project (pom.xml or build.gradle found)
2. Run tests: `mvn test` or `gradle test`
3. Run coverage: `mvn verify` (JaCoCo) or `gradle test jacocoTestReport`
4. Run linting: `mvn checkstyle:check` or `gradle checkstyleMain`
5. AI validators check architecture, errors, complexity

### 4. Review Quality

```bash
@review feature-auth
```

SDP will run AI validators:
- `/coverage-validator` - Analyzes JaCoCo coverage
- `/architecture-validator` - Checks Clean Architecture layers
- `/error-validator` - Finds empty catch blocks
- `/complexity-validator` - Identifies complex methods

## Quality Gates

### Test Coverage ≥80%

**Maven:**
```bash
mvn verify
```

**Gradle:**
```bash
gradle test jacocoTestReport
```

**Output:**
```
[INFO] --- jacoco-maven-plugin:0.8.11:check (default) @ myproject ---
[INFO] Loading execution data file target/jacoco.exec
[INFO] Analyzed bundle 'myproject' with 150 classes
[INFO] All coverage checks have been met.
[INFO]   Rule 0: CoveredRatio = 0.85 (required minimum = 0.80) ✅ PASS
```

### Type Checking

**Maven:**
```bash
mvn compiler:compile
```

**Gradle:**
```bash
gradle compileJava
```

**Output:**
```
[INFO] --- maven-compiler-plugin:3.11.0:compile (1 sources) ---
[INFO] Nothing to compile - all classes are up to date
✅ PASS
```

### Linting

**Maven:**
```bash
mvn checkstyle:check
```

**Gradle:**
```bash
gradle checkstyleMain
```

**Output:**
```
[INFO] --- checkstyle:6.3.1:check (default) @ myproject ---
[INFO] Starting audit...
[INFO] Audit done.
✅ PASS
```

### File Size

```bash
find src/main -name "*.java" -exec wc -l {} + | awk '$1 > 200'
```

**Output:**
```
(no output = all files <200 LOC) ✅ PASS
```

## Example Workflow

```bash
# Start feature
@feature "Add user login"

# Plan workstreams
@design feature-login

# Execute first workstream
@build 00-001-01

# Expected output:
# ✓ Project type detected: Java (pom.xml)
# ✓ Running tests: mvn test
# ✓ Coverage: 85% (≥80% required)
# ✓ Type checking: javac -Xlint:all
# ✓ Linting: mvn checkstyle:check
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!

# Execute next workstream
@build 00-001-02

# Review all workstreams
@review feature-login

# Deploy
@deploy feature-login
```

## Maven Configuration

**pom.xml additions:**
```xml
<build>
  <plugins>
    <!-- JaCoCo for coverage -->
    <plugin>
      <groupId>org.jacoco</groupId>
      <artifactId>jacoco-maven-plugin</artifactId>
      <version>0.8.11</version>
      <executions>
        <execution>
          <goals>
            <goal>prepare-agent</goal>
          </goals>
        </execution>
        <execution>
          <id>report</id>
          <phase>test</phase>
          <goals>
            <goal>report</goal>
          </goals>
        </execution>
        <execution>
          <id>check</id>
          <goals>
            <goal>check</goal>
          </goals>
          <configuration>
            <rules>
              <rule>
                <element>PACKAGE</element>
                <limits>
                  <limit>
                    <counter>LINE</counter>
                    <value>COVEREDRATIO</value>
                    <minimum>0.80</minimum>
                  </limit>
                </limits>
              </rule>
            </rules>
          </configuration>
        </execution>
      </executions>
    </plugin>

    <!-- Checkstyle for linting -->
    <plugin>
      <groupId>org.apache.maven.plugins</groupId>
      <artifactId>maven-checkstyle-plugin</artifactId>
      <version>3.3.0</version>
      <configuration>
        <configLocation>google_checks.xml</configLocation>
      </configuration>
    </plugin>
  </plugins>
</build>
```

## Gradle Configuration

**build.gradle additions:**
```groovy
plugins {
    id 'jacoco'
}

test {
    useJUnitPlatform()
    finalizedBy jacocoTestReport
}

jacoco {
    toolVersion = "0.8.11"
}

jacocoTestReport {
    dependsOn test
    reports {
        xml.required = true
        csv.required = false
        html.outputLocation = layout.buildDirectory.dir('reports/jacoco')
    }
}

jacocoTestCoverageVerification {
    violationRules {
        rule {
            limit {
                minimum = 0.80
            }
        }
    }
}

tasks.named('check') {
    dependsOn jacocoTestCoverageVerification
}
```

## Common Issues

### Issue: Coverage <80%

**Solution:** Add more tests
```java
// src/test/java/com/example/service/UserServiceTest.java
@Test
public void testCreateUser() {
    User user = userService.create("Alice");
    assertEquals("Alice", user.getName());
}

@Test
public void testCreateUserInvalidName() {
    assertThrows(IllegalArgumentException.class,
        () -> userService.create(""));  // Add this test
}
```

### Issue: Empty catch blocks

**Solution:** Log or re-throw
```java
// Before (FAIL)
try {
    riskyOperation();
} catch (Exception e) {
    // Empty catch
}

// After (PASS)
try {
    riskyOperation();
} catch (ValueException e) {
    logger.error("Invalid value", e);
    throw e;
}
```

### Issue: Checkstyle violations

**Solution:** Fix style issues
```bash
mvn checkstyle:checkstyle
# Review violations and fix
```

## Tips

1. **Use Maven wrapper:**
   ```bash
   ./mvnw test
   ```

2. **Run specific test:**
   ```bash
   mvn test -Dtest=UserServiceTest
   ```

3. **Skip coverage during development:**
   ```bash
   mvn test -DskipTests=false -Djacoco.skip=true
   ```

4. **Generate coverage HTML report:**
   ```bash
   mvn jacoco:report
   open target/site/jacoco/index.html
   ```

## Next Steps

- [Python Quick Start](../python/QUICKSTART.md)
- [Go Quick Start](../go/QUICKSTART.md)
- [Full Tutorial](../../TUTORIAL.md)
