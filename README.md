<div align="center">

<h1>mutabaah-cli</h1>
<p><strong>CLI sederhana untuk mencatat aktivitas ibadah dan rutinitas harian.</strong></p>
<p>Go · JSON lokal · Tanpa dependency eksternal</p>
<p>
  <img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
  <img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />
</p>

</div>

---

*Mutabaah* adalah kegiatan mencatat dan mengevaluasi amalan. Proyek ini merupakan mini project untuk mempelajari Go melalui aplikasi CLI yang mengelola data nyata: argumen command line, struct dan slice, JSON persistence, target sesi, dan validasi input.

## Fitur

- Menambahkan aktivitas dengan ID otomatis dan target yang dapat diatur per aktivitas.
- Mencatat sesi aktivitas melalui `catat <id>`; progress bertambah satu pada setiap pencatatan.
- Menandai target selesai otomatis ketika jumlah sesi tercapai, atau secara manual dengan `sudah <id>`.
- Melihat aktivitas belum selesai dengan `list`, yang sudah selesai dengan `list --done`, atau semuanya dengan `list --all`.
- Menampilkan indikator progress, target sesi, dan peringatan `TERLAMBAT` jika target waktu lewat dan aktivitas belum selesai.
- Mengubah judul dan menghapus aktivitas.
- Menyimpan data di JSON lokal dan memindahkan file lama ke lokasi penyimpanan baru saat diperlukan.

## Persyaratan

- [Go](https://go.dev/dl/) **1.26.5 atau lebih baru**, sesuai deklarasi pada `go.mod`.
- Hanya menggunakan standard library Go; tidak ada dependency eksternal.

## Mulai

Clone repositori dan masuk ke folder proyek:

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

Di Windows, binary dapat dibuat sebagai `mutabaah-cli.exe`:

```powershell
go build -o mutabaah-cli.exe .
.\mutabaah-cli.exe help
```

## Perintah

Jalankan perintah dari direktori proyek jika menggunakan `go run .`, atau gunakan path binary yang sudah dibuat.

| Perintah | Fungsi |
|---|---|
| `tambah "judul aktivitas" [--sesi N] [--periode HARI]` | Membuat aktivitas dengan target sesi dan periode |
| `list` | Menampilkan aktivitas yang belum selesai |
| `list --done` | Menampilkan aktivitas yang sudah selesai |
| `list --all` | Menampilkan semua aktivitas |
| `catat <id>` | Mencatat satu sesi progress |
| `sudah <id>` | Menuntaskan target secara manual |
| `ubah <id> "judul baru"` | Mengubah judul aktivitas |
| `hapus <id>` | Menghapus aktivitas |
| `help` / `-h` / `--help` | Menampilkan bantuan tanpa membaca atau mengubah file data |

Alias `list -all` dan `list -done` juga didukung. Flag `--all` dan `--done` tidak dapat digunakan bersamaan.

### Target aktivitas

Target default adalah **3 sesi dalam 7 hari**, dihitung sejak aktivitas dibuat. Atur target khusus saat menambah aktivitas:

```bash
go run . tambah "Tilawah satu juz" --sesi 5 --periode 10
```

- `--sesi` adalah jumlah pencatatan yang dibutuhkan untuk menyelesaikan target.
- `--periode` adalah tenggat dalam hari sejak aktivitas dibuat.
- `catat <id>` menaikkan progress satu sesi. Ketika progress mencapai target sesi, aktivitas otomatis ditandai selesai.
- `sudah <id>` adalah shortcut manual: aktivitas langsung selesai dan progress disamakan dengan target sesi.
- Jika tenggat lewat sementara aktivitas belum selesai, `list` menampilkan `TERLAMBAT`. Pencatatan tetap diperbolehkan setelah tenggat; periodenya tidak di-reset otomatis.

### Contoh penggunaan

```bash
$ go run . tambah "Tilawah satu juz" --sesi 3 --periode 7
Aktivitas "Tilawah satu juz" dengan ID 1 berhasil ditambahkan (target 3 sesi dalam 7 hari)

$ go run . tambah "Sholat dhuha"
Aktivitas "Sholat dhuha" dengan ID 2 berhasil ditambahkan (target 3 sesi dalam 7 hari)

$ go run . catat 1
Sesi dicatat untuk aktivitas 1. Progress: 1/3 sesi

$ go run . catat 1
Sesi dicatat untuk aktivitas 1. Progress: 2/3 sesi

$ go run . catat 1
Sesi dicatat untuk aktivitas 1. Progress: 3/3 sesi
Target sesi tercapai; aktivitas ditandai selesai.

$ go run . list
|Ncan| 2. Sholat dhuha [0/3 sesi; target 7 hari]

$ go run . list --done
|Ntos| 1. Tilawah satu juz [3/3 sesi; target 7 hari]

$ go run . list --all
|Ntos| 1. Tilawah satu juz [3/3 sesi; target 7 hari]
|Ncan| 2. Sholat dhuha [0/3 sesi; target 7 hari]
```

Status yang ditampilkan:

- `|Ncan|` — aktivitas belum selesai.
- `|Ntos|` — aktivitas sudah selesai.

Aktivitas yang selesai disembunyikan dari tampilan default, bukan dihapus. Gunakan `list --all` untuk melihatnya kembali.

## Penyimpanan data

Data disimpan pada direktori home pengguna:

```text
~/.mtcli/daftaraktivitas.json
```

Direktori dan file dibuat otomatis saat perubahan data perlu disimpan. Command `help`, `-h`, dan `--help` berhenti sebelum penentuan lokasi data, migrasi, atau pembacaan JSON. Karena itu, bantuan tetap dapat ditampilkan meskipun file data rusak.

### Migrasi penyimpanan lama

Jika lokasi penyimpanan baru belum ada dan aplikasi menemukan `daftaraktivitas.json` di direktori kerja saat ini, file tersebut disalin ke lokasi baru di folder home. File lama tidak dihapus otomatis.

### Bentuk data

```json
{
  "daftarAktivitas": [
    {
      "id": 1,
      "judul": "Tilawah satu juz",
      "target": {
        "sesi": 3,
        "periode": 7
      },
      "proses": 2,
      "ketuntasan": false,
      "dibuatSaat": "2026-10-08T07:30:00+09:00"
    }
  ]
}
```

Field `target.sesi` dan `target.periode` diserialisasi dengan nama lowercase. File lama yang menggunakan `Sesi` dan `Periode` masih dapat dibaca; saat data berikutnya disimpan, key ditulis dengan format lowercase. Aktivitas lama yang sudah berstatus selesai tetapi belum memiliki progress akan dinormalisasi agar progress sesuai target.

## Struktur kode

| File | Tanggung jawab |
|---|---|
| `main.go` | Routing perintah, parsing argumen, output CLI, dan penyimpanan setelah perubahan |
| `mutabaah.go` | Model aktivitas, target, progress, status selesai, tenggat, dan operasi daftar |
| `inventory.go` | Lokasi file, pembacaan/penulisan JSON, normalisasi data lama, dan migrasi |
| `main_test.go` | Pengujian parsing target, flag list, dan kompatibilitas JSON |
| `mutabaah_test.go` | Pengujian progress sesi dan penyelesaian target |

## Catatan desain

- Perubahan pada daftar dilakukan melalui pointer receiver `ListMutabaah`.
- Progress dibatasi maksimum pada jumlah sesi target; pencatatan tidak bisa menambah progress setelah aktivitas selesai.
- `sudah <id>` tetap tersedia sebagai penyelesaian manual untuk alur kerja yang tidak memerlukan pencatatan setiap sesi.
- ID baru dihitung dari ID aktivitas terakhir yang masih ada ditambah satu. Karena itu, menghapus aktivitas terakhir dapat menyebabkan ID tersebut digunakan kembali.
- Periode merupakan tenggat satu kali dari tanggal aktivitas dibuat, bukan periode berulang yang otomatis melakukan reset.

## Lisensi

Dirilis di bawah [MIT License](LICENSE).
