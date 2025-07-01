package decrypt
import ( "database/sql"; "fmt"; _ "github.com/mutecomm/go-sqlcipher/v4" )
func DumpDB(src, dst string, key []byte) error {
 dsn := fmt.Sprintf("file:%s?_pragma_key=x'%x'&_pragma_cipher_page_size=4096&_pragma_kdf_iter=64000", src, key)
 db, err := sql.Open("sqlite3", dsn); if err != nil { return err }; defer db.Close()
 _, err = db.Exec(fmt.Sprintf("ATTACH DATABASE '%s' AS dec KEY ''", dst))
 if err != nil { return err }
 _, err = db.Exec("SELECT sqlcipher_export('dec')")
 return err
}