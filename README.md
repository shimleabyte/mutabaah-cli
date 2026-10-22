<div align="center">

# mutabaah-cli

**CLI sederhana untuk mencatat aktivitas ibadah dan rutinitas harian.**

Dibuat dengan Go · Penyimpanan lokal berbasis JSON · Tanpa dependency eksternal

<img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
<img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />

</div>

---

*Mutabaah* adalah kegiatan mencatat dan mengevaluasi amalan. Proyek ini merupakan mini project untuk mempelajari Go melalui aplikasi CLI yang benar-benar mengelola data: membaca argumen, memanipulasi slice dan struct, menyimpan JSON, serta menangani lokasi file dan migrasi data lama.

## Fitur

- Menambahkan aktivitas dengan ID otomatis.
- Menampilkan aktivitas yang belum selesai.
- Menampilkan semua aktivitas dengan `list --all`.
- Menandai aktivitas sebagai selesai.
- Mengubah judul dan menghapus aktivitas.
- Menyimpan data secara lokal dalam JSON yang mudah dibaca.
- Memindahkan data dari file lama `daftaraktivitas.json` ke lokasi penyimpanan baru saat diperlukan.

## Persyaratan

- [Go](https://go.dev/dl/) **1.26.5 atau lebih baru**, sesuai deklarasi versi pada `go.mod`.

Proyek ini hanya menggunakan standard library Go; tidak ada dependency eksternal yang perlu diunduh.

## Mulai

Clone repositori dan masuk ke direktori proyek:

```bash
git clone https://github.com/shimleabyte/mutabaah-cli.git
cd mutabaah-cli
```

Jalankan langsung tanpa membuat binary:

```bash
go run . help
```

Atau build binary:

```bash
go build -o mutabaah-cli .
./mutabaah-cli help
```

Di Windows, nama binary hasil build dapat menggunakan ekstensi `.exe`:

```powershell
go build -o mutabaah-cli.exe .
.\mutabaah-cli.exe help
```

## Perintah

Jalankan perintah dari direktori proyek ketika menggunakan `go run .`, atau gunakan path binary yang sudah dibuat.

| Perintah | Fungsi |
|---|---|
| `tambah "judul aktivitas"` | Membuat aktivitas baru |
| `list` | Menampilkan aktivitas yang belum selesai |
| `list --all` | Menampilkan aktivitas selesai dan belum selesai |
| `sudah <id>` | Menandai aktivitas sebagai selesai |
| `ubah <id> "judul baru"` | Mengubah judul aktivitas |
| `hapus <id>` | Menghapus aktivitas |
| `help` | Menampilkan bantuan |

Alias `list -all` juga didukung.

## Contoh penggunaan

```bash
$ go run . tambah "Tilawah satu juz"
Aktivitas Tilawah satu juz dengan ID 1 berhasil ditambahkan

$ go run . tambah "Sholat dhuha"
Aktivitas Sholat dhuha dengan ID 2 berhasil ditambahkan

$ go run . list
|Ncan| 1. Tilawah satu juz
|Ncan| 2. Sholat dhuha

$ go run . sudah 1
Bagoes! Aktivitas 1 sudah ditandai terlaksana

$ go run . list
|Ncan| 2. Sholat dhuha

$ go run . list --all
|Ntos| 1. Tilawah satu juz
|Ncan| 2. Sholat dhuha

$ go run . ubah 2 "Sholat dhuha 4 rakaat"
Aktivitas 2 berhasil diubah menjadi "Sholat dhuha 4 rakaat"

$ go run . hapus 2
Aktivitas 2 dihapus
```

Status yang ditampilkan oleh `list`:

- `|Ncan|` — aktivitas belum selesai.
- `|Ntos|` — aktivitas sudah selesai.

Aktivitas yang selesai disembunyikan dari tampilan default, tetapi tidak dihapus. Gunakan `list --all` untuk melihatnya kembali.

## Penyimpanan data

Data disimpan pada direktori home pengguna, di dalam:

```text
~/.mtcli/daftaraktivitas.json
```

Di Windows, lokasi ini berada di direktori home pengguna, pada folder `.mtcli\daftaraktivitas.json`.

Direktori dan file dibuat otomatis ketika ada perubahan data yang perlu disimpan. Perintah yang hanya membaca data seperti `list` tidak menulis ulang file JSON.

### Migrasi dari penyimpanan lama

Jika file penyimpanan baru belum ada dan aplikasi menemukan `daftaraktivitas.json` di direktori kerja saat ini, file tersebut akan disalin ke lokasi baru di folder home. File lama tidak dihapus otomatis.

### Bentuk data

Contoh struktur JSON yang digunakan:

```json
{
  "daftarAktivitas": [
    {
      "id": 1,
      "judul": "Tilawah satu juz",
      "target": {
        "Sesi": 3,
        "Periode": 7
      },
      "proses": 0,
      "ketuntasan": false,
      "dibuatSaat": "2026-10-08T07:30:00+09:00"
    }
  ]
}
```

## Struktur kode

| File | Tanggung jawab |
|---|---|
| `main.go` | Entry point, routing perintah, validasi argumen, dan penyimpanan setelah perubahan |
| `mutabaah.go` | Model `Aktivitas` dan `ListMutabaah`, serta operasi tambah, selesai, ubah, dan hapus |
| `inventory.go` | Penentuan lokasi data, pembacaan/penulisan JSON, dan migrasi file lama |

## Catatan implementasi

- Operasi perubahan data dilakukan melalui pointer receiver pada `ListMutabaah`, sehingga perubahan diterapkan pada slice yang dikelola.
- Data hanya disimpan kembali setelah perintah yang mengubah daftar berhasil dijalankan.
- ID baru dihitung dari ID aktivitas terakhir yang masih ada ditambah satu. Menghapus aktivitas terakhir dapat membuat ID tersebut digunakan kembali.
- Perintah saat ini hanya dapat menandai aktivitas sebagai selesai; belum tersedia perintah untuk membatalkan status selesai.
- Field `target` dan `proses` sudah disimpan dalam JSON, tetapi belum dikelola oleh perintah CLI.

## Lisensi

Dirilis di bawah [MIT License](LICENSE).
