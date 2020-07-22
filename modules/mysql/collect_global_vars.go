package mysql

import (
	"strconv"
)

const (
	queryGlobalVariables = "SHOW GLOBAL VARIABLES"
)

/*
MariaDB [(none)]> SHOW GLOBAL VARIABLES;
+------------------+-------+
| Variable_name    | Value |
+------------------+-------+
| max_connections  | 151   |
| table_open_cache | 2000  |
+------------------+-------+
*/

var globalVariablesMetrics = []string{
	"max_connections",
	"table_open_cache",
}

func (m *MySQL) collectGlobalVariables(collected map[string]int64) error {
	rows, err := m.db.Query(queryGlobalVariables)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range globalVariablesMetrics {
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalVariable(strValue)
		if err != nil {
			continue
		}
		collected[name] = value
	}
	return nil
}

func parseGlobalVariable(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}