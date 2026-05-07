# pgwatch copilot

**pgwatch copilot** is a command-line assistant for exploring PostgreSQL monitoring data from a pgwatch sink.  
It helps you inspect metrics, understand system behavior, and ask natural-language questions about your monitored databases, while keeping the sink connection read-only for safety.

## Features

- Connects to a pgwatch sink database
- Enforces read-only access at the PostgreSQL level
- Retrieves available `sys_id` values from the sink
- Uses an LLM to help explain monitoring data
- Supports configurable lookback windows
- Designed to help with query and schema context for monitored PostgreSQL instances
