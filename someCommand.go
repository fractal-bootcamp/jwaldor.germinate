// import (
//     "io/ioutil"
//     "os"
//     "path/filepath"
// )

// func copyDir(src string, dst string) error {
//     entries, err := ioutil.ReadDir(src)
//     if err != nil {
//         return err
//     }
//     for _, entry := range entries {
//         srcPath := filepath.Join(src, entry.Name())
//         dstPath := filepath.Join(dst, entry.Name())

//         if entry.IsDir() {
//             if err := os.MkdirAll(dstPath, 0755); err != nil {
//                 return err
//             }
//             if err := copyDir(srcPath, dstPath); err != nil {
//                 return err
//             }
//         } else {
//             if err := copyFile(srcPath, dstPath); err != nil {
//                 return err
//             }
//         }
//     }
//     return nil
// }

// func copyFile(src, dst string) error {
//     data, err := ioutil.ReadFile(src)
//     if err != nil {
//         return err
//     }
//     return ioutil.WriteFile(dst, data, 0644)
// }

// func createDatabaseBoilerplate() error {
//     src := "assets/database-boilerplate-main/"
//     dst := "path/to/output/directory"
//     return copyDir(src, dst)
// }