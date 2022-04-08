# Pembuat
- Nama : Dipta Laksmana Baswara
- NPM : 1806235832

# Hal yang dibutuhkan
- docker-compose

# Cara Menjalankan
1. Mulai server postgresql
```shell
docker-compose -f docker-compose.yml up -d postgres
```
2. Mulai server read dan write
```shell
docker-compose -f docker-compose.yml up -d server
```
3. Mulai server nginx
```shell
docker-compose -f docker-compose.yml up -d nginx
```

Catatan: Semua server harus dimulai secara berurutan

# Analisis Perbandingan
- Screenshot uji coba Idempoten API
![Idempoten](Idempoten%20Test.png)
- Screenshot uji coba Non Idempoten API
![Non Idempoten](Non%20Idempoten%20Test.png)
  
Dari hasil diatas dapat diambil kesimpulan bahwa API Idempoten yang menggunakan cache 
lebih bisa menerima request lebih banyak (hampir 3x lipat) daripada API Non Idempoten 
yang tidak menggunakan API, hal tersebut terjadi karena pada API Idempoten, cache 
digunakan sebagai tempat penyimpanan sementara untuk data yang kemungkinan sering 
dibutuhkan oleh client dan disimpan pada memory yang berkecepatan tinggi, sehingga 
pada pemanggilan yang sama sebelum cache tersebut expired nginx akan mengambil data 
pada cache untuk mempersingkat waktu.  
