
package decrypt

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DecryptDB uses the provided key to decrypt an encrypted SQLite database
// and saves the decrypted version to dstPath.
func DecryptDB(key, srcPath, dstPath string) error {
	// Ensure the source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source database not found: %s", srcPath)
	}

	// Remove destination file if it exists, to ensure a clean state
	_ = os.Remove(dstPath)

	// 1. Open the encrypted database with the key
	srcDSN := fmt.Sprintf("file:%s?_pragma_key=%s&_pragma_cipher_page_size=4096", srcPath, key)
	db, err := sql.Open("sqlite3", srcDSN)
	if err != nil {
		return fmt.Errorf("failed to open encrypted database: %w", err)
	}
	defer db.Close()

	// 2. Verify the key by trying to access data
	if _, err := db.Exec("SELECT count(*) FROM sqlite_master;"); err != nil {
		return fmt.Errorf("failed to verify key, likely incorrect: %w", err)
	}

	// 3. Attach a new, unencrypted database
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS decrypted KEY ''", dstPath)
	if _, err := db.Exec(attachSQL); err != nil {
		return fmt.Errorf("failed to attach decrypted database: %w", err)
	}

	// 4. Export the content from the encrypted DB to the unencrypted one
	if _, err := db.Exec("SELECT sqlcipher_export('decrypted');"); err != nil {
		return fmt.Errorf("failed to export database content: %w", err)
	}

	// 5. Detach the new database
	if _, err := db.Exec("DETACH DATABASE decrypted;"); err != nil {
		return fmt.Errorf("failed to detach decrypted database: %w", err)
	}

	return nil
}
