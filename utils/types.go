package utils

// Definisikan tipe khusus untuk Jenis Aktivitas
type JenisAktivitas string

// Opsi pilihan jenis aktivitas yang valid (Enum)
const (
	JenisChecklist JenisAktivitas = "checklist" // Tinggal centang (0 atau 1)
	JenisCounter   JenisAktivitas = "counter"   // Input angka progres (misal: minum 5 gelas)
	JenisNote      JenisAktivitas = "note"      // Berupa keterangan/catatan teks
)

type Activity struct {
	Nama    string         `json:"nama"`
	Jenis   JenisAktivitas `json:"jenis"`
	Target  int            `json:"target"`  // Jika Checklist -> Target = 1. Jika Note -> Target = 0 atau diabaikan.
	Progres int            `json:"progres"` // Default otomatis 0 di Go saat struct dibuat
	Catatan string         `json:"catatan"` // Tambahan field jika Jenisnya adalah "note"
}
