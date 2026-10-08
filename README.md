# mutabaah-cli

Tracker ibadah/aktivitas harian sederhana berbasis command line, ditulis dengan Go (standard library saja, tanpa dependency eksternal). Data disimpan lokal dalam file JSON.

*Mutabaah* (muhasabah/mutaba'ah) = mencatat dan mengevaluasi amalan harian. Project ini dibuat sebagai latihan belajar Go: struct, pointer receiver, slice, JSON persistence, dan CLI routing via `os.Args`.

## Fitur

- Tambah aktivitas baru
- Lihat semua aktivitas beserta statusnya
- Tandai aktivitas sudah dilaksanakan
- Ubah judul aktivitas
- Hapus aktivitas
- Data tersimpan otomatis ke file JSON (human-readable)

## Prasyarat

- [Go](https://go.dev/dl/) 1.21+ (disarankan versi terbaru)

## Instalasi

```bash
git clone https://github.com/<username>/mutabaah-cli.git
cd mutabaah-cli
go build -o mutabaah-cli .
```

Atau langsung jalankan tanpa build:

```bash
go run . list
```

## Cara Pakai

```text
mutabaah-cli tambah "judul aktivitas"    Tambah aktivitas baru
mutabaah-cli list                        Tampilkan semua aktivitas
mutabaah-cli sudah <id>                  Tandai aktivitas selesai
mutabaah-cli ubah <id> "judul baru"      Ubah judul aktivitas
mutabaah-cli hapus <id>                  Hapus aktivitas
mutabaah-cli help                        Tampilkan bantuan
```

### Contoh

```bash
$ ./mutabaah-cli tambah "Tilawah satu juz"
Aktivitas "Tilawah satu juz" berhasil ditambahkan dengan ID 1

$ ./mutabaah-cli tambah "Sholat dhuha"
Aktivitas "Sholat dhuha" berhasil ditambahkan dengan ID 2

$ ./mutabaah-cli list
|Ncan| 1. Tilawah satu juz
|Ncan| 2. Sholat dhuha

$ ./mutabaah-cli sudah 1
Bagoes! Aktivitas 1 sudah ditandai terlaksana

$ ./mutabaah-cli list
|Ntos| 1. Tilawah satu juz
|Ncan| 2. Sholat dhuha

$ ./mutabaah-cli ubah 2 "Sholat dhuha 4 rakaat"
Aktivitas 2 berhasil diubah menjadi "Sholat dhuha 4 rakaat"

$ ./mutabaah-cli hapus 2
Aktivitas 2 dihapus
```

### Arti status

| Status   | Arti                      |
|----------|---------------------------|
| `\|Ncan\|` | Belum dilaksanakan (*can* = belum) |
| `\|Ntos\|` | Sudah dilaksanakan (*ntos* = sudah) |

## Penyimpanan Data

Data disimpan di file `daftaraktivitas.json` pada **direktori tempat perintah dijalankan**. Kalau file belum ada, akan dibuat otomatis saat perintah pertama yang mengubah data dijalankan.

Contoh isi file:

```json
{
  "daftarAktivitas": [
    {
      "id": 1,
      "judul": "Tilawah satu juz",
      "target": { "Sesi": 3, "Periode": 7 },
      "proses": 0,
      "ketuntasan": true,
      "dibuatSaat": "2026-10-08T07:30:00+09:00"
    }
  ]
}
```

> Karena path-nya relatif, menjalankan CLI dari folder berbeda akan membaca/menulis file yang berbeda pula.

## Struktur Project

| File           | Isi                                                                 |
|----------------|---------------------------------------------------------------------|
| `main.go`      | Entry point, routing command via `os.Args`, parsing ID, teks bantuan |
| `mutabaah.go`  | Struct `Aktivitas`, `UkurTarget`, `ListMutabaah` + method `Tambah`, `Sudah`, `Hapus`, `Ubah` |
| `inventory.go` | Persistensi JSON: `LoadAktivitas` dan `SimpanAktivitas`             |

## Catatan Desain

- **Pointer receiver** pada `ListMutabaah` supaya method memodifikasi data asli, bukan salinan.
- **ID auto-increment**: ID baru = ID aktivitas terakhir + 1, jadi ID yang sudah dihapus tidak dipakai ulang selama aktivitas terakhir masih ada.
- **Hanya menyimpan jika ada perubahan**: command read-only seperti `list` dan `help` tidak menulis ulang file.
- Field `target` (`Sesi`, `Periode`) dan `proses` sudah ada di struct dan JSON, tapi belum dipakai oleh command apa pun (disiapkan untuk pengembangan berikutnya).

## Lisensi

MIT License