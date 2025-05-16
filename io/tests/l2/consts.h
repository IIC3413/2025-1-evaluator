#pragma once

#include "relational_model/schema.h"
#include <cstdint>
#include <string>
#include <vector>

constexpr uint64_t GB = 1024 * 1024 * 1024;
constexpr uint64_t BUFF_SIZE = 1 * GB;
const std::string DB_DIR = "data/eval_dbs";
const std::string INPUTS_DIR = "inputs/eval_dbs";
const std::string OUTPUT_DIR = "outputs";
const std::vector<ColumnInfo> COLUMNS = {
    {"name", DataType::STR},
    {"level", DataType::INT},
    {"class", DataType::STR},
};
