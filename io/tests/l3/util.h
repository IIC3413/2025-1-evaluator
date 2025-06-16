#pragma once

#include "query/parser/parser.h"
#include "query/translator/optimizer.h"
#include "relational_model/schema.h"
#include "storage/heap_file/heap_file.h"
#include "system/system.h"
#include "tests/consts.h"
#include <algorithm>
#include <filesystem>
#include <functional>
#include <string>
#include <vector>

namespace Utils {
static System init_system(std::string& test_name) {
  return System::init(std::string(DB_DIR + "/" + test_name), BUFF_SIZE);
}

[[maybe_unused]]
static HeapFile* create_table(std::string table_name, std::vector<ColumnInfo>&& columns) {
  Schema schema(std::move(columns));
  Schema existing_table_schema;
  HeapFile* table = catalog.get_table(table_name, &existing_table_schema);
  if (table == nullptr) {
    table = catalog.create_table(table_name, schema);
  } else {
    assert(existing_table_schema == schema);
  }
  return table;
}

static std::string record_ref_to_string(RecordRef& ref) {
  std::string ref_str;
  for (auto it = ref.values.begin(); it != ref.values.end(); ++it) {
    auto val = *it.base();

    if (val->datatype == DataType::INT) {
      ref_str += std::to_string(val->value.as_int);
    } else {
      ref_str += "\"";
      ref_str += val->value.as_str;
      ref_str += "\"";
    }

    if (it != ref.values.end()) {
      ref_str += ",";
    }
  }
  return ref_str;
}

static void run_query(std::string& test_name, std::string& query) {
  if (!std::filesystem::exists(OUTPUT_DIR)) {
    std::filesystem::create_directory(OUTPUT_DIR);
  }

  auto logical_plan = Parser::parse(query, false);
  auto query_iter = Optimizer::create_physical_plan(std::move(logical_plan));

  // Header
  std::string header_string;
  auto out_cols = query_iter->get_columns();
  if (out_cols.size() > 0) {
    header_string += out_cols[0].alias + '.' + out_cols[0].info.name;
  }
  for (size_t i = 1; i < out_cols.size(); i++) {
    header_string += ',' + out_cols[i].alias + '.' + out_cols[i].info.name;
  }

  // Write output
  auto& out = query_iter->get_output();
  auto out_strings = std::vector<std::string>();
  query_iter->begin();
  while (query_iter->next()) {
    out_strings.push_back(record_ref_to_string(out));
  }
  std::sort(out_strings.begin(), out_strings.end());

  std::ofstream ofstream(OUTPUT_DIR + "/" + test_name + ".output");
  ofstream << header_string << '\n';
  for (const auto& str : out_strings) {
    ofstream << str << '\n';
  }
}

[[maybe_unused]]
static void run(std::string& test_name, std::function<void()> fn, std::string& query) {
  auto system = init_system(test_name);
  fn();
  run_query(test_name, query);
}
} // namespace Utils
