# F068-F073 UX Protocol Checklist

This checklist ensures consistent UX improvements across all F068-F073 features.

## Pre-Implementation Checklist

- [ ] Read UX baseline report (`docs/reference/2026-02-16-f068-ux-baseline.md`)
- [ ] Identify affected friction points
- [ ] Define measurable success criteria
- [ ] Plan measurement integration

## During Implementation

- [ ] Follow clean architecture principles
- [ ] Maintain test coverage >= 80%
- [ ] Keep files under 200 LOC
- [ ] Add inline help text for new commands
- [ ] Include examples in command documentation

## Post-Implementation Verification

- [ ] TTFV measurement possible via evidence log
- [ ] Error messages include fix commands
- [ ] Help text is discoverable and consistent
- [ ] Status output points to next actions
- [ ] Demo/template walkthrough succeeds

## Release Gate Criteria

### F068 (UX Foundation)
- [ ] `sdp init --guided` succeeds on clean environment
- [ ] `sdp --help` shows commands grouped by intent
- [ ] Template walkthrough completes in < 5 minutes
- [ ] All UX metrics within target ranges

### F069 (Error Messaging)
- [ ] All error messages include fix commands
- [ ] Error resolution time < 3 minutes
- [ ] Recovery success rate > 90%

### F070 (Progress Feedback)
- [ ] Long operations show progress
- [ ] Cancellation is graceful
- [ ] Status is always available

### F071 (Documentation Sync)
- [ ] CLI help and docs share vocabulary
- [ ] All commands documented with examples
- [ ] Getting-started guide is verified

### F072 (Accessibility)
- [ ] Color is not the only indicator
- [ ] Screen reader compatible
- [ ] Keyboard navigation works

### F073 (Internationalization)
- [ ] Message extraction complete
- [ ] Locale detection works
- [ ] Date/time formatting localized

## Metrics Tracking

| Feature | TTFV Target | Discoverability Target | Status |
|---------|-------------|----------------------|--------|
| F068 | < 15 min | 90% | In Progress |
| F069 | < 12 min | 92% | Pending |
| F070 | < 15 min | 90% | Pending |
| F071 | < 10 min | 95% | Pending |
| F072 | < 15 min | 90% | Pending |
| F073 | < 15 min | 90% | Pending |

---

**Version:** 1.0
**Last Updated:** 2026-02-16
