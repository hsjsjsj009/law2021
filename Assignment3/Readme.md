# Cara Menjalankan
1. Jalankan docker-compose, harus menggunakan docker-compose karena saya perlu menggunakan jaringan docker
```bash
docker-compose -f docker-compose.yml up
```
2. Jalankan main.go, hal ini harus di lakukan pertama kali, karena exchange baru dibuat di main.go
```bash
go run main.go
```
3. Buka halaman localhost:5000 di browser yang diinginkan

# Penjelasan Tambahan
Saya disini menggunakan reverse proxy di nginx karena terdapat masalah CORS yang disebabkan oleh perbedaan host yakni disini websocket berada pada localhost:15674 sedangkan web server di localhost:5000