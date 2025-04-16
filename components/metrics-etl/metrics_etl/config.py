"""
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
"""
import os

AGENT_ID = os.getenv("AGENT_ID", "default_agent")
ZENOH_CONFIG_FILE = os.getenv("ZENOH_CONFIG_FILE", "config/zenoh_config.json5")
VERSION = os.getenv("VERSION", "develop")
PROMETHEUS_PORT = int(os.getenv("PROMETHEUS_PORT", 8999))
