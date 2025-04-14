#pragma once

#include "relational_model/schema.h"
#include "storage/filesystem.h"
#include "storage/heap_file/heap_file.h"
#include "storage/heap_file/heap_file_iter.h"
#include "system/system.h"
#include <cstdint>
#include <filesystem>
#include <functional>
#include <string>
#include <utility>
#include <vector>

constexpr uint64_t GB = 1024 * 1024 * 1024;
constexpr uint64_t BUFF_SIZE = 1 * GB;
const std::string DB_DIR = "data/eval_dbs";
const std::string OUTPUT_DIR = "outputs";
const std::vector<ColumnInfo> COLUMNS = {
    {"name", DataType::STR},
    {"level", DataType::INT},
    {"class", DataType::STR},
};

namespace Utils {

static std::string get_db_name(std::string test_name) {
  return std::string(DB_DIR + "/" + test_name);
}

static HeapFile* get_or_create_table(std::string table_name, Schema& schema) {
  Schema existing_table_schema;
  HeapFile* table = catalog.get_table(table_name, &existing_table_schema);

  if (table == nullptr) {
    table = catalog.create_table(table_name, schema);
  } else {
    assert(existing_table_schema == schema);
  }

  return table;
}

static HeapFile*
clone_table(std::string test_name, const Schema& schema, std::string src_name, std::string target_name) {
  Schema src_schema;
  assert(catalog.get_table(src_name, &src_schema) != nullptr);
  std::filesystem::copy_file(
      std::filesystem::path(get_db_name(test_name) + "/" + src_name),
      std::filesystem::path(get_db_name(test_name) + "/" + target_name),
      std::filesystem::copy_options::update_existing
  );
  return catalog.create_table(target_name, schema);
}

static std::pair<std::string, std::string> test_table_names(std::string& test_name) {
  auto input_name = test_name + std::string("_input");
  auto output_name = test_name + std::string("_output");
  return std::pair(input_name, output_name);
}

[[maybe_unused]]
static void zero_out_pages_free_space(std::string& table_name) {
  Schema schema;
  auto heap_file = catalog.get_table(table_name, &schema);
  uint64_t total_pages = file_mgr.count_pages(heap_file->file_id);
  char* zeros = new char[Page::SIZE]();
  for (uint64_t i = 0; i < total_pages; i++) {
    HeapFilePage page(heap_file->file_id, i);
    auto entries = page.get_dir_count();
    auto free_space = page.get_free_space();
    page.page.write((2 + entries) * sizeof(int32_t), free_space, zeros);
  }
  delete []zeros;
}

[[maybe_unused]]
static void generate_test_tables(
    std::string& test_name,
    Schema& schema,
    std::function<void(std::string&)> input_setup,
    std::function<void(std::string&)> output_setup
) {
  // Initialize the system and set variables.
  auto [input_table_name, output_table_name] = test_table_names(test_name);
  get_or_create_table(input_table_name, schema);
  // Setup the input table that submissions will run on.
  input_setup(input_table_name);
  // Write to disk.
  buffer_mgr.flush();
  // Create the expected output table.
  clone_table(test_name, schema, input_table_name, output_table_name);
  // Modify the expected output table to its desired state.
  output_setup(output_table_name);
}

[[maybe_unused]]
static void run_test(std::string& test_name, std::function<void(std::string&)> run) {
  // Run operations over input file.
  auto [input_table_name, output_table_name] = test_table_names(test_name);
  run(input_table_name);
  // Write to disk.
  buffer_mgr.flush();
  // Clone resulting table to outuputs.
  if (!Filesystem::is_directory(OUTPUT_DIR)) {
    Filesystem::create_directories(OUTPUT_DIR);
  }
  std::filesystem::copy_file(
      std::filesystem::path(get_db_name(test_name) + "/" + input_table_name),
      std::filesystem::path(OUTPUT_DIR + "/" + output_table_name),
      std::filesystem::copy_options::update_existing
  );
}
} // namespace Utils
