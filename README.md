<div align="center">

<h1>mutabaah</h1>
<p><strong>CLI untuk mencatat aktivitas harian dengan target sesi dan tenggat.</strong></p>
<p>Go · Penyimpanan JSON lokal · Tanpa dependency eksternal</p>
<p>
  <img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
  <img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />
</p>

</div>

---

*Mutabaah* adalah CLI kecil untuk mencatat aktivitas harian. Perintah dipisahkan berdasarkan domain: `acts` untuk pengelolaan aktivitas dan `notes` sebagai ruang terpisah untuk fitur catatan yang direncanakan.

Proyek ini juga menjadi sarana belajar Go melalui implementasi praktis: routing command, struct dan slice, persistensi JSON, validasi input, pencatatan progress, serta migrasi data.

## Fitur

- Menambahkan aktivitas dengan target sesi dan tenggat yang bisa diatur.
- Mencatat sesi yang dilakukan dan menampilkan progress bar di terminal.
- Menandai aktivitas selesai otomatis ketika target sesi tercapai.
- Menampilkan aktivitas yang belum selesai, sudah selesai, atau semuanya.
- Menampilkan aktivitas yang melewati tenggat.
- Mengubah judul dan menghapus aktivitas.
- Menyimpan data secara lokal dalam JSON serta memigrasikan file data lama jika diperlukan.
- Menampilkan bantuan tanpa membaca atau mengubah data aktivitas.

## Persyaratan

- [Go](https://go.dev/dl/) **1.26.5 atau lebih baru**, sesuai deklarasi di `go.mod`.
- Hanya menggunakan standard library Go; tidak ada dependency eksternal.

## Memulai

Clone repositori:

```bash
git clone https://github.com/shimleabyte/mutabaah-cli.git
cd mutabaah-cli
```

Jalankan tanpa membuat binary:

```bash
go run . help
```

Atau build executable dengan nama yang lebih pendek, `mutabaah`:

```bash
go build -o mutabaah .
./mutabaah help
```

Di Windows:

```powershell
go build -o mutabaah.exe .
.\\mutabaah.exe help
```

## Struktur command

Command root membagi fitur berdasarkan domain.

| Command | Fungsi |
|---|---|
| `mutabaah acts list` | Menampilkan aktivitas |
| `mutabaah acts tambah "Judul aktivitas"` | Menambahkan aktivitas |
| `mutabaah acts catat <id>` | Mencatat satu sesi |
| `mutabaah acts sudah <id>` | Menandai aktivitas selesai secara manual |
| `mutabaah acts ubah <id> "Judul baru"` | Mengubah judul aktivitas |
| `mutabaah acts hapus <id>` | Menghapus aktivitas |
| `mutabaah acts help` | Menampilkan bantuan domain aktivitas |
| `mutabaah notes` | Placeholder domain catatan; belum tersedia |
| `mutabaah help` | Menampilkan bantuan root |

Domain `notes` disiapkan sebagai bagian terpisah. Saat ini, command tersebut belum membuat, membaca, mengubah, atau menghapus catatan.

## Target aktivitas

Target default aktivitas baru adalah **3 sesi dalam 7 hari**. Target per aktivitas dapat ditetapkan menggunakan `--sesi` dan `--periode`:

```bash
./mutabaah acts tambah "Tilawah satu juz" --sesi 5 --periode 10
```

- `--sesi` adalah jumlah sesi yang diperlukan untuk menyelesaikan target.
- `--periode` adalah tenggat satu kali dalam hari, dihitung sejak aktivitas dibuat.
- `catat <id>` menambah progress satu sesi.
- Ketika progress mencapai target, aktivitas otomatis ditandai selesai.
- `sudah <id>` tetap tersedia untuk menyelesaikan aktivitas secara manual.
- Tenggat tidak otomatis di-reset. Periode berulang belum diimplementasikan.

## Menampilkan aktivitas

Secara default, list hanya menampilkan aktivitas yang belum selesai.

```bash
./mutabaah acts list
./mutabaah acts list --done
./mutabaah acts list --all
```

Contoh output:

```text
ID  STATE  ACTIVITY           PROGRESS          DEADLINE
1   OPEN   Tilawah satu juz   [####------] 2/5  5d left
2   LATE   Murajaah           [##--------] 1/5  late 3d
3   DONE   Sholat Dhuha       [##########] 3/3  complete
```

| Status | Arti |
|---|---|
| `OPEN` | Aktivitas belum selesai dan tenggat belum lewat |
| `LATE` | Aktivitas belum selesai dan tenggat sudah lewat |
| `DONE` | Target sesi tercapai atau aktivitas diselesaikan manual |

Progress bar menggunakan `#` untuk sesi yang tercatat dan `-` untuk sisa sesi. `--done` hanya menampilkan aktivitas selesai, sedangkan `--all` menampilkan semuanya.

## Penyimpanan data

Data aktivitas disimpan di direktori home pengguna:

```text
~/.mtcli/daftaraktivitas.json
```

Direktori dan file dibuat ketika ada perubahan yang perlu disimpan. Command `mutabaah help` dan `mutabaah acts help` berhenti sebelum mencari lokasi data, menjalankan migrasi, atau membaca JSON. Jadi, bantuan tetap tersedia ketika file JSON rusak.

### Migrasi data lama

Jika file penyimpanan baru belum ada dan `daftaraktivitas.json` ditemukan di direktori kerja saat ini, file lama akan disalin ke lokasi baru. File asal tidak dihapus secara otomatis.

### Format JSON

```json
{
  "daftarAktivitas": [
    {
      "id": 1,
      "judul": "Tilawah satu juz",
      "target": {
        "sesi": 5,
        "periode": 10
      },
      "proses": 2,
      "ketuntasan": false,
      "dibuatSaat": "2026-10-23T07:30:00Z"
    }
  ]
}
```

Key target menggunakan huruf kecil (`sesi` dan `periode`). File lama dengan key `Sesi` dan `Periode` tetap dapat dibaca.

## Struktur proyek

| File | Tanggung jawab |
|---|---|
| `main.go` | Routing root/domain, parsing argumen, output tabel, dan command aktivitas |
| `mutabaah.go` | Model aktivitas, progress sesi, penyelesaian target, dan tenggat |
| `inventory.go` | Lokasi data, persistensi JSON, normalisasi data lama, dan migrasi |
| `main_test.go` | Tes parsing argumen, kompatibilitas JSON, dan routing bantuan |
| `mutabaah_test.go` | Tes progress sesi, penyelesaian, dan tenggat |
| `table_test.go` | Tes progress bar, output tabel, dan routing command |

## Batasan desain saat ini

- Domain `acts` sudah tersedia; `notes` masih berupa placeholder.
- Periode merupakan tenggat satu kali, bukan siklus berulang.
- ID baru dihitung dari ID aktivitas terakhir yang masih ada, sehingga menghapus aktivitas terakhir dapat membuat ID tersebut dipakai kembali.
- `sudah <id>` dapat menuntaskan aktivitas tanpa mencatat setiap sesi satu per satu.

## Lisensi

Dirilis di bawah [MIT License](LICENSE).
