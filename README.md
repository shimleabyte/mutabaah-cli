# mutabaah-cli

### This is my first 'learning by doing' mini Go project. 

### Project structure:
mutabaah-cli/
├── main.go
├── activity.go
├── stats.go
├── tracker.go
└── types.go 

This is mutabaah-cli, daily prayer/habit tracker command line interface built in go language.
"Mutabaah" is just different terminology of habit tracker specialized in islamic prayers, so it's more like daily prayer tracker, often used by islamic education culture in indonesia.

### Project purposes -> Implement how go declaration works:
variable declaration
type declaration
array/slice
pointer
function type
function variable
closure 

### Feature:
Saving Activities
Adding
Summarize Total Activities
Using Function as Data/State

Current commit status as journal:
# Go Learning Project — Declaration Syntax

## Tujuan

Project ini dibuat sebagai media belajar Go secara bertahap. Setiap materi Go langsung diterapkan ke dalam project kecil, sehingga setiap konsep memiliki alasan keberadaan dan tidak hanya dipelajari melalui contoh yang terpisah.

Project yang dibangun adalah **Mutaba'ah CLI Tracker**, namun domain dibuat sesederhana mungkin agar fokus utama tetap pada konsep bahasa Go.

---

# Materi 1 — Declaration Syntax

Materi ini mengacu pada artikel Rob Pike **"Go's Declaration Syntax"**.

## Filosofi Go

Go menggunakan pendekatan deklarasi yang dibaca **dari kiri ke kanan**, mengikuti cara manusia membaca.

Contoh:

```go
var age int
```

Dibaca sebagai:

> Variabel `age` bertipe `int`.

Begitu juga untuk fungsi:

```go
func Add(a int, b int) int
```

Dibaca sebagai:

> Function `Add` menerima dua `int` dan mengembalikan sebuah `int`.

Berbeda dengan C yang mendeskripsikan tipe melalui bentuk ekspresi, Go memisahkan sintaks tipe dan sintaks ekspresi sehingga deklarasi menjadi lebih mudah dibaca ketika tipe mulai kompleks.

---

# Penerapan ke Project

## 1. Struct Declaration

Membuat model data pertama.

```go
type Activity struct {
    Name     string
    Target   int
    Progress int
}
```

Objek ini merepresentasikan sebuah aktivitas yang memiliki:

* Nama aktivitas
* Target
* Progress saat ini

---

## 2. Struct Initialization

Membuat instance pertama.

```go
var act1 Activity = Activity{
    Name:     "Tadarus",
    Target:   5,
    Progress: 3,
}
```

Pemahaman penting:

* `Activity` adalah **type (blueprint)**.
* `act1` adalah **instance (object/value)**.

---

## 3. Printing Struct

Menggunakan package `fmt`.

```go
fmt.Printf("%+v\n", act1)
```

Output:

```text
{Name:Tadarus Target:5 Progress:3}
```

Formatter yang dipelajari:

* `%v` → hanya nilai.
* `%+v` → nama field + nilai.
* `%#v` → representasi lengkap dalam sintaks Go.

---

## 4. Function Declaration

Membuat fungsi yang menerima sebuah `Activity`.

Contoh konsep:

```text
CompletionPercentage(activity Activity) int
```

Dibaca sebagai:

> Function menerima sebuah `Activity` dan mengembalikan sebuah `int`.

---

## 5. Accessing Struct Fields

Field diakses melalui **instance**, bukan melalui **type**.

Benar:

```go
activity.Progress
activity.Target
```

Salah:

```go
Activity.Progress
```

Karena `Activity` adalah nama tipe, bukan objek yang menyimpan data.

---

## 6. Return Value

Menghitung persentase progress.

Contoh logika:

```text
percentage = (Progress × 100) / Target
```

Nilai tersebut kemudian dikembalikan menggunakan `return`.

---

## 7. Runtime Error: Division by Zero

Operasi pembagian integer dengan penyebut nol akan menyebabkan panic.

Contoh:

```go
5 / 0
0 / 0
```

Menghasilkan runtime error:

```text
panic: runtime error: integer divide by zero
```

Solusinya adalah melakukan validasi sebelum pembagian.

Contoh aturan bisnis sederhana:

* Jika `Target == 0`, kembalikan `0`.

Perlu dibedakan antara:

* **Matematika** → pembagian dengan nol tidak terdefinisi.
* **Program** → operasi tersebut menyebabkan runtime error sehingga harus ditangani oleh programmer.

---

# Konsep yang Dipelajari

* Package
* Struct
* Exported field
* Variable declaration
* Struct literal
* Function declaration
* Function parameter
* Return value
* Field access
* Integer arithmetic
* Basic validation
* fmt.Printf formatting

---

# Struktur Project Saat Ini

```text
mutabaah-cli/

├── main.go
├── types.go
├── tracker.go
└── stats.go
```

---

# Selanjutnya

Materi berikutnya:

* Method
* Receiver
* Memindahkan perilaku dari fungsi menjadi method milik `Activity`

Transformasi yang akan dilakukan:

```text
Sebelum

CompletionPercentage(activity)

↓

Sesudah

activity.CompletionPercentage()
```

Ini merupakan langkah awal menuju pemrograman yang lebih berorientasi objek dengan pendekatan idiomatis Go.


