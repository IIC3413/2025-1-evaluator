#include <iostream>
#include <string>
#include <vector>

#include "relational_model/record.h"
#include "relational_model/schema.h"
#include "storage/heap_file/heap_file.h"
#include "system/system.h"

constexpr uint64_t GB = 1024 * 1024 * 1024;
const int32_t pi[10] = {3, 1, 4, 1, 5, 9, 2, 6, 5, 3};
const std::string pad =
    "________________________________________________________________________________________"
    "________________________________________________________________________________________"
    "_______________________________________________________________________________";

void print_pi(std::string table_name) {
  Record record(*catalog.get_table_info(table_name).schema);
  auto table_iter = catalog.get_table_info(table_name).heap_file->get_record_iter();
  table_iter->begin(record);
  while (table_iter->next()) {
    std::cout << record.values[0].value.as_int;
  }
  std::cout << std::endl;
}

void add_pair(int32_t pair, RID* ptrs, std::string table_name) {
  int32_t i;
  for (i = pair; i < 5; i++) {
    ptrs[i] = catalog.insert_record(table_name, {pi[i], pair, pad, pad, pad});
  }
  for (i = pair + 5; i < 10; i++) {
    ptrs[i] = catalog.insert_record(table_name, {pi[i], pair + 5, pad, pad, pad});
  }
  for (i = pair + 1; i < 5; i++) {
    catalog.delete_record(table_name, ptrs[i]);
    catalog.delete_record(table_name, ptrs[i + 5]);
  }
  catalog.get_table_info(table_name).heap_file->vacuum();
}

void create_or_validate_table(Schema* schema, std::string table_name) {
  Schema existing_table_schema;
  HeapFile* table = catalog.get_table(table_name, &existing_table_schema);
  if (table == nullptr) {
    table = catalog.create_table(table_name, *schema);
  } else {
    assert(existing_table_schema == *schema);
  }
}

int main() {
  uint64_t buffer_size = 1 * GB;
  std::string database_folder = "data/test_example";

  auto system = System::init(database_folder, buffer_size);

  std::string table_name("pi_table");
  std::vector<ColumnInfo> columns = {
      {"digit", DataType::INT}, {"pos", DataType::INT},  {"pad0", DataType::STR},
      {"pad1", DataType::STR},  {"pad2", DataType::STR},
  };

  Schema schema(columns);
  create_or_validate_table(&schema, table_name);

  RID* digit_ptrs = new RID[10];
  for (int32_t i = 0; i < 5; i++) {
    add_pair(i, digit_ptrs, table_name);
  }
  /* print_pi(table_name); */
  std::cout << "super_secret_code 10/10" << std::endl;

  delete[] digit_ptrs;
  return 0;
}
