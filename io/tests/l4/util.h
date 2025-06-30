#pragma once

#include "query/optimizer/optimizer.h"
#include "query/parser/parser.h"
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
  std::ofstream ofstream(OUTPUT_DIR + "/" + test_name + ".output");

  auto logical_plan = Parser::parse(query, false);
  auto query_iter = Optimizer::create_physical_plan(std::move(logical_plan));

  // Header
  auto out_cols = query_iter->get_columns();
  if (out_cols.size() > 0) {
    ofstream << out_cols[0].alias << '.' << out_cols[0].info.name;
  }
  for (size_t i = 1; i < out_cols.size(); i++) {
    ofstream << ',' << out_cols[i].alias << '.' << out_cols[i].info.name;
  }
  ofstream << '\n';

  // Write output
  auto& out = query_iter->get_output();
  auto out_strings = std::vector<std::string>();
  query_iter->begin();
  while (query_iter->next()) {
    out_strings.push_back(record_ref_to_string(out));
  }
  std::sort(out_strings.begin(), out_strings.end());
  for (const auto& str : out_strings) {
    ofstream << str << "\n";
  }
}

[[maybe_unused]]
static void run(std::string& test_name, std::function<void()> fn, std::string& query) {
  auto system = init_system(test_name);
  fn();
  run_query(test_name, query);
}

[[maybe_unused]]
static void run_bloom_filter(std::string& test_name, std::function<void()> fn, std::string& query) {
  if (!std::filesystem::exists(OUTPUT_DIR)) {
    std::filesystem::create_directory(OUTPUT_DIR);
  }
  std::ofstream ofstream(OUTPUT_DIR + "/" + test_name + ".output");
  std::cout.rdbuf(ofstream.rdbuf());

  auto system = init_system(test_name);
  fn();

  auto logical_plan = Parser::parse(query, false);
  auto query_iter = Optimizer::create_physical_plan(std::move(logical_plan));
  auto out_strings = std::vector<std::string>();
  query_iter->begin();
  while (query_iter->next()) {
  }
}
} // namespace Utils
