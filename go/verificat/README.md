# Verificat

## About

- **Meaning**: Short for "Verify Catalog"
- **Pronunciation**: rhymes with "Terrific Cat"
- **Functionality**: This app is an autonomous agent that performs independent continuous verification tests against the catalog of Weedmaps Services running in production environments. This set of tests is known as a **Production Readiness Checklist**.

## Concept

The idea behind Verificat is to build a system that does two things:

1. Know a desired State of Readiness for a Weedmaps Service, and
2. Report on the reality of that State of Readiness.

### State of Readiness

When a Service is considered "Production Ready" it embodies eight quantifiable (that is, measurable) principles. **Acting together** these drive the availability of a Service. Each of these Eight Principles have about so many subdivisions that can be enumerated.

From the top level, they are:

* *stability*
* *reliability*
* *scalability*
* *performance*
* *fault tolerance*
* *catastrophe-preparedness*
* *monitoring*
* *documentation*

Any test performed by Verificat will fall into one or more of these categories.

The outcome we want to measure is our **State of Readiness**. To be done at any time, continuously.

Susan Fowler describes how we need to think about Readiness:

> The basic idea behind production-readiness is this: a production-ready application or service is one that can be trusted to serve production traffic. When we refer to an application or microservice as “production-ready,” we confer a great deal of trust upon it: we trust it to behave reasonably, we trust it to perform reliably, we trust it to get the job done and to do its job well with very little downtime.

### Scoring

On each run, the Service being tested starts with a fresh score of 100. The full suite of tests are run and a Score is derived by subtracting 1 from the total for each failed verification. Successful verifications leave the Score as-is.

Higher scores have more coverage, but the goal isn't to enforce a Score of 100. Instead, we want to show a Service can continuously display its State of Readiness.

The service being tested will receive a new score each time a request to test is triggered. For this reason, the only visible score in the database is the most recent. In future versions we want to add the ability to keep a timeseries database of run IDs and scores.

#### Required Baseline

In the future we want to tag tests as "Required" in order to build a minimum baseline for scoring to "allow" a Service to either be "freshly deployed" as a checklist, or to grade a Service to "remain in Production" but with the penalty of maintenace to raise the Score.

#### Why 100?

Think of the Score as a _grade_, not as a _percent completed_.

The grading of a Service can be quite complex! A *single test* can be shown to cover *multiple points* along the Eight Principles. Instead of aggregating test results that will be skewed by the number of tests, our golf game style scoring allows us to have as little as one test and still see a meaningful Score.

The number is also our upper bound on test count. If we need more than 100 tests to verify 64 quantifiable outcomes, we're probably doing math wrong or we need to reconsider the boundary of the system being tested.

### Objectives

Eventually we will build an overall "service readiness objective" that all Services must pass. There may be several (or many) "optional" tests that boost one of the Eight Principles. Deficiencies in some tests may be counter-balanced by increased coverage, for instance.

#### Current WIP: Owner

The only item being tested is the equality of the Owner field in Backstage. Verificat uses the source of truth for this value, the GitHub CODEOWNERS file, for comparison.

### Test-Driven Development

Because this tool is a Test-Driven approach to Production Checklists, the development approach to building this automation tool is also Test-Driven (TDD).

TDD allows us to be iterative. While building the automation, we need to make sure we're doing exactly what we want, to both run the test and be confident in the accuracy of the data.

Our initial approach is to tackle one test (or type) at a time, but Verificat will run tests in parallel.


## Verify Catalog

### Verification vs Validation

1. ***Validation*** is for displaying that a thing meets the **design requirements** of the system.
2. ***Verification*** seeks evidence that a thing meets the **outcome requirements** of the system.

A good reference on the sometimes subtle difference, the [V&V Wikipedia article](https://en.wikipedia.org/wiki/Verification_and_validation) provides some other examples. This automation reports both Validation and Verification, for instance:

1. The contents of the Owner field is _validated_ as being present. Not checked for *correctness*, only that the field contains a non-null value.
2. The *correctness* of the Owner field is _verified_ with an independent check against some source of truth, however that may present itself (e.g.: data lookup, run a function, even initiate a process like chaos engineering), and report the measurement outcome.

## Autonomy

Verificat seeks to be as independent as possible so that it can measure as closely to real-world as possible. For this reason it is meant to be run as an autonomous service that can perform any number of actions against real-world infrastructure.

It is not meant to be dependent on any other service, nor provide dependencies downstream. It is a data gatherer, analyst, and presenter. Its inputs are the Services we run and the tests we want to display a State of Readiness. The outputs are the Scores.

### Notifications

A facility to send email alerting would be a good initial notification service.

## Observability

Verificat uses the Golang stdlib `slog`. Plans are to use the OpenTelemetry Handler extension to emit telemetry. This is the choice in order to make the code as non-proprietary as possible, furthermore to support CNCF and Open Source.

Metrics won't be as important as Logs with Traces.

## Operations

1. Run locally with: `docker run -ti --rm --name verificat -p 4330:4330 verificat:latest`
2. In another terminal, run a test against the `admin` service: `curl -X POST http://localhost:4330/v0/admin`
4. Get results for all services: `curl http://localhost:4330/v0/almanac`

## References

- [Wiki Page](https://weedmaps.atlassian.net/wiki/spaces/SRE/pages/31868059667/Production+Readiness+Checklist) in the SRE Confluence Space
