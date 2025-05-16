#pragma once

#include "relational_model/schema.h"
#include "storage/b_plus_tree/b_plus_tree.h"
#include "storage/b_plus_tree/b_plus_tree_dir.h"
#include "storage/b_plus_tree/b_plus_tree_leaf.h"
#include "storage/heap_file/heap_file.h"
#include "system/system.h"
#include "tests/consts.h"
#include <cstdint>
#include <filesystem>
#include <fstream>
#include <functional>
#include <ios>
#include <string>

using test_fn = std::function<void(std::string&)>;

namespace Utils {
static System init_system(std::string& test_name) {
  return System::init(std::string(DB_DIR + "/" + test_name), BUFF_SIZE);
}

static HeapFile* init_db(std::string& table_name, Schema& schema) {
  Schema existing_table_schema;
  HeapFile* table = catalog.get_table(table_name, &existing_table_schema);
  if (table == nullptr) {
    table = catalog.create_table(table_name, schema);
  } else {
    assert(existing_table_schema == schema);
  }
  return table;
}

static void zero_out_index_dirs(FileId id) {
  char* zeros = new char[Page::SIZE]();
  auto& end_page = buffer_mgr.append_page(id);
  end_page.unpin();
  for (int32_t i = 0; i < end_page.page_id.page_number; i++) {
    auto& page = buffer_mgr.get_page(id, static_cast<int64_t>(i));
    auto cc = page.read_int32(0);
    auto child_offset = BPlusTreeDir::OFFSET_CHILDREN + cc * sizeof(int32_t);
    auto record_offset = BPlusTreeDir::OFFSET_RECORD + (cc - 1) * sizeof(BPlusTreeRecord);
    page.write(child_offset, BPlusTreeDir::OFFSET_RECORD - child_offset, zeros);
    page.write(record_offset, Page::SIZE - record_offset, zeros);
    page.unpin();
  }
  delete []zeros;
}

static void zero_out_index_leafs(FileId id) {
  char* zeros = new char[Page::SIZE]();
  auto& end_page = buffer_mgr.append_page(id);
  end_page.unpin();
  for (int32_t i = 0; i < end_page.page_id.page_number; i++) {
    auto& page = buffer_mgr.get_page(id, static_cast<int64_t>(i));
    auto rc = page.read_int32(0);
    auto record_offset = BPlusTreeLeaf::OFFSET_RECORDS + rc * sizeof(BPlusTreeRecord);
    page.write(record_offset, Page::SIZE - record_offset, zeros);
    page.unpin();
  }
  delete []zeros;
}

static void zero_out_index(std::string& test_name) {
  auto& table_info =  catalog.get_table_info(test_name);
  auto bpt = dynamic_cast<BPlusTree*>(table_info.index.get());
  if (bpt == nullptr) {
    return;
  }
  zero_out_index_dirs(bpt->dir_file_id);
  zero_out_index_leafs(bpt->leaf_file_id);
}

[[maybe_unused]]
static void clone_db(std::string& test_name) {
  auto origin = DB_DIR + "/" + test_name + "/";
  auto target = INPUTS_DIR + "/" + test_name + "/";
  // If the inputs directory does not exist create it.
  if (!std::filesystem::is_directory(INPUTS_DIR)) {
    std::filesystem::create_directories(INPUTS_DIR);
  }
  std::filesystem::copy(
      std::filesystem::path(origin), std::filesystem::path(target),
      std::filesystem::copy_options::update_existing | std::filesystem::copy_options::recursive
  );
}

[[maybe_unused]]
static void clone_output(std::string& test_name) {
  if (!std::filesystem::is_directory(OUTPUT_DIR)) {
    std::filesystem::create_directories(OUTPUT_DIR);
  }
  std::ifstream if_bpt_dir(DB_DIR + "/" + test_name + "/" + test_name + ".bpt.dir", std::ios_base::binary);
  std::ifstream if_bpt_leaf(DB_DIR + "/" + test_name + "/" + test_name + ".bpt.leaf", std::ios_base::binary);
  std::ofstream of_concat(OUTPUT_DIR + "/" + test_name + ".bpt.output", std::ios_base::binary);
  of_concat << if_bpt_dir.rdbuf() << if_bpt_leaf.rdbuf();
}

[[maybe_unused]]
static void run(std::string& test_name, test_fn fn, test_fn post) {
  {
    Schema schema(COLUMNS);
    auto system = init_system(test_name);
    init_db(test_name, schema);
    fn(test_name);
    zero_out_index(test_name);
  }
  post(test_name);
}
} // namespace Utils
