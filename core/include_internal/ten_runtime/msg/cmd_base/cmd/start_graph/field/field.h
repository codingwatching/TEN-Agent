//
// Copyright © 2025 Agora
// This file is part of TEN Framework, an open source project.
// Licensed under the Apache License, Version 2.0, with certain conditions.
// Refer to the "LICENSE" file in the root directory for more information.
//
#pragma once

#include "ten_runtime/ten_config.h"

typedef enum TEN_CMD_START_GRAPH_FIELD {
  TEN_CMD_START_GRAPH_FIELD_CMD_HDR,

  TEN_CMD_START_GRAPH_FIELD_LONG_RUNNING_MODE,
  TEN_CMD_START_GRAPH_FIELD_EXTENSION_INFO,

  TEN_CMD_START_GRAPH_FIELD_PREDEFINED_GRAPH,
  TEN_CMD_START_GRAPH_FIELD_GRAPH_JSON,

  TEN_CMD_START_GRAPH_FIELD_LAST,
} TEN_CMD_START_GRAPH_FIELD;
