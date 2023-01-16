# Go-Attendance

## Instructions

1. Install `Go` dengan versi `1.19`.

2. Install package yang digunakan pada aplikasi ini.
```go
go mod download
```

3. Buka MySQL editor dan jalankan script yang ada di direktori `database/migration` dengan file yang bernama `db.sql`.

4. Build aplikasi ini.
```go
go build
```

5. Jalankan aplikasi.
```bash
./go-attendance.exe
```

6. Aplikasi dapat diakses melalui `http://localhost:8000`

## Dokumentasi API Aplikasi Go-Attendance

### 1.  Register Karyawan

#### Explanation

Register karyawan dilakukan untuk membuat akun baru karyawan agar karyawan tersebut dapat login ke dalam sistem absensi ini. Diperlukan email yang belum pernah terdaftar dalam sistem agar akun dapat dibuat dimana pada implementasi pada UserService akan dicek apakah terdapat email yang sama. Jika tidak maka akun akan dibuat dengan cara memasukkan data ke dalam table Users.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/register
```
#### Body

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| email   | required	  	| email yang bersifat unique dari karyawan |
| full_name         | required      |  nama lengkap karyawan |
| password   | required	  	| password yang ingin diberikan |

#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| id dari user yang berhasil ditambahkan |
|email| email dari user yang berhasil ditambahkan |
|full_name| nama lengkap dari user yang berhasil ditambahkan|

#### Example
```json
curl https://go-attendance.herokuapp.com/register -X POST -d '{"email": "sam@gmail.com","password": "samuel","full_name": "sam"}' -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 1,
        "email": "sam@gmail.com",
        "full_name": "sam"
    }
}
```

### 2.  Login Karyawan

#### Explanation

Login karyawan dilakukan agar karyawan dapat masuk kedalam akunnya dan dapat mengakses fitur-fitur yang ada di sistem ini. Nantinya password akan dihash dan dibandingkan dengan password yang mempunyai email yang sama dengan email yang diinput karyawan. Lalu jika ternyata password cocok, maka sistem akan membuat session cookie yang berlaku selama 1 hari dimana session tersebut berisi JWT authentication. Dengan adanya cookie, maka sistem dapat mengenal bahwa karyawan masih login kedalam akunnya, sehingga karyawan tidak perlu memasukkan email dan passwordnya berulang kali. Kenapa menggunakan JWT authentication, karena pada metode autentikasi ini dapat menyimpan informasi sederhana dari user yaitu email dan user ID sehingga nantinya memudahkan agar tidak perlu mengambil email dan user ID berulang kali dari database.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/login
```
#### Body

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| email   | required	  	| email dari karyawan |
| password   | required	  	| password yang dimasukkan |

#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| id dari user yang berhasil ditambahkan |
|message| pesan yang memberitahukan bahwa user berhasil login |

#### Example
```json
curl https://go-attendance.herokuapp.com/login -X POST -d '{"email": "sam@gmail.com","password": "samuel"}' -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 1,
        "message": "successfully login to account"
    }
}
```

### 3.  Logout Karyawan

#### Explanation

Logout karyawan berguna untuk menghapus session karyawan sehingga karyawan tersebut harus memasukkan email dan passwordnya kembali. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya maka session akan dihapus dari cookie.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/logout
```


#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|message| pesan yang memberitahukan user berhasil logout dari akun|

#### Example
```json
curl https://go-attendance.herokuapp.com/logout -X POST -H "Content-Type: application/json"
```
```json
{
    "data": {
        "message": "successfully logout from account"
    }
}
```

### 4.  Check In Karyawan

#### Explanation

Check-In karyawan berguna untuk memasukkan data jam masuk karyawan pada hari tersebut. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya akan dicek lagi dimana terdapat middleware yang berguna untuk mengecek apakah karyawan sudah melakukan check-in pada hari tersebut. Jika karyawan belum melakukan check-in maka data check in yang baru untuk hari tersebut akan dibuat.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/attendance/checkin
```

#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data check in yang berhasil dibuat|
|attendance_date| tanggal kehadiran user pada hari tersebut dengan format `YYYY-MM-dd`|
|check_in_time| jam masuk user pada hari tersebut dengan format `HH:mm:ss`|
|message| pesan yang ditampilkan|


#### Example
```json
curl https://go-attendance.herokuapp.com/attendance/checkin -X POST -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 3,
        "attendance_date": "2023-01-16",
        "message": "successfully check in today",
        "check_in_time": "07:27:15"
    }
}
```

### 5.  Check Out Karyawan

#### Explanation

Check-Out karyawan berguna untuk memasukkan data jam keluar karyawan pada hari tersebut. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya akan dicek lagi dimana terdapat middleware yang berguna untuk mengecek apakah karyawan sudah melakukan check-in pada hari tersebut. Jika karyawan belum melakukan check-in maka endpoint ini akan mengirimkan pesan bahwa karyawan harus check-in terlebi dahulu. Sebaliknya, karyawan akan dicek apakah telah melakukan check-out untuk memastikan karyawan tidak melakukan check-out lebih dari sekali. Kemudian data check-out akan dipasangkan dengan data check-in pada hari tersebut di database.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/attendance/checkout
```

#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data check in yang berhasil dibuat|
|attendance_date| tanggal kehadiran user pada hari tersebut `YYYY-MM-dd`|
|check_out_time| jam keluar user pada hari tersebut `HH:mm:ss`|
|message| pesan yang ditampilkan|


#### Example
```json
curl https://go-attendance.herokuapp.com/attendance/checkout -X POST -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 3,
        "attendance_date": "2023-01-16",
        "message": "successfully check in today",
        "check_out_time": "07:32:15"
    }
}
```

### 6.  List Kehadiran Karyawan

#### Explanation

List Kehadiran karyawan berguna untuk mendapatkan semua data kehadiran karyawan yang tercatat dalam sistem. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya maka data data kehadiran karyawan akan diambil dan ditampilkan.

#### HTTP Request
```json
GET https://go-attendance.herokuapp.com/attendance
```

#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data kehadiran yang berhasil dibuat|
|attendance_date| tanggal kehadiran user pada hari tersebut `YYYY-MM-dd`|
|check_out_time| jam keluar user pada hari tersebut `HH:mm:ss`|
|check_in_time| jam masuk user pada hari tersebut `HH:mm:ss`|
|message| pesan yang ditampilkan|
|total_attendances| jumlah kehadiran dari user tersebut|


#### Example
```json
curl https://go-attendance.herokuapp.com/attendance -X GET -H "Content-Type: application/json"
```
```json
{
    "data": {
        "attendances": [
            {
                "id": 4,
                "attendance_date": "2023-01-16",
                "check_out_time": "07:35:25",
                "check_in_time": "07:35:19"
            }
        ],
        "total_attendances": 1
    }
}
```

### 7.  Insert Aktivitas Karyawan

#### Explanation

Insert Aktivitas Karyawan berguna untuk memasukkan aktivitas-aktivitas karyawan yang dilakukan pada hari tersebut. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya akan dicek lagi dimana terdapat middleware yang berguna untuk mengecek apakah karyawan telah check-in pada hari tersebut. Jika belum maka karyawan akan disarankan untuk melakukan check-in, namun jika sudah maka karyawan dapat menambahkan aktivitas pada hari tersebut.

#### HTTP Request
```json
POST https://go-attendance.herokuapp.com/activity
```

#### Body

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| description   | required	  	| deskripsi aktivitas yang dilakukan pada hari tersebut |


#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data aktivitas yang berhasil dibuat|
|attendance_id| `id` dari data kehadiran pada hari tersebut |
|description| deskripsi aktivitas yang dilakukan pada hari tersebut|

#### Example
```json
curl https://go-attendance.herokuapp.com/activity -X POST -d '{"description": "Pada hari ini memperbaiki bug yang ada"}' -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 1,
        "attendance_id": 1,
        "description": "Pada hari ini memperbaiki bug yang ada"
    }
}
```

### 8.  Update Aktivitas Karyawan

#### Explanation

Update Aktivitas Karyawan berguna untuk mengubah data aktivitas karyawan. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya akan dicek lagi dimana terdapat middleware yang berguna untuk mengecek apakah karyawan telah check-in pada hari tersebut. Jika belum maka karyawan akan disarankan untuk melakukan check-in, namun jika sudah maka karyawan dapat mengubah data aktivitas yang ingin dirubah.

#### HTTP Request
```json
PATCH https://go-attendance.herokuapp.com/activity/{id}
```

#### Parameters

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| id   | required	  	| `id` dari data aktivitas yang ingin diubah |


#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data aktivitas yang berhasil diubah|
|description| deskripsi aktivitas yang dilakukan pada hari tersebut|

#### Example
```json
curl https://go-attendance.herokuapp.com/activity/1 -X PATCH -d '{"description": "Pada hari ini, bugnya sudah selesai diperbaiki"}' -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 1,
        "description": "Pada hari ini, bugnya sudah selesai diperbaiki"
    }
}
```

### 9.  Delete Aktivitas Karyawan

#### Explanation

Delete Aktivitas Karyawan berguna untuk menghapus data aktivitas karyawan. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya akan dicek lagi dimana terdapat middleware yang berguna untuk mengecek apakah karyawan telah check-in pada hari tersebut. Jika belum maka karyawan akan disarankan untuk melakukan check-in, namun jika sudah maka karyawan dapat menghapus data aktivitas yang ingin dihapus.

#### HTTP Request
```json
DELETE https://go-attendance.herokuapp.com/activity/{id}
```

#### Parameters

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| id   | required	  	| `id` dari data aktivitas yang ingin dihapus |


#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data aktivitas yang berhasil dihapus|

#### Example
```json
curl https://go-attendance.herokuapp.com/activity/1 -X DELETE -H "Content-Type: application/json"
```
```json
{
    "data": {
        "id": 3
    }
}
```

### 10.  List Aktivitas Karyawan Berdasarkan Tanggal Kehadiran

#### Explanation

List Aktivitas Karyawan Berdasarkan Tanggal Kehadiran berguna untuk mendapatkan list data aktivitas karyawan berdasarkan tanggal kehadiran karyawan. Namun ketika mengakses endpoint ini, terdapat middleware yang mengecek apakah karyawan sudah login atau belum. Jika karyawan belum login namun mengakses endpoint ini maka karyawan akan diberitahu bahwa karyawan belum login, sebaliknya tanggal kehadiran yang dimasukkan akan dicari id dari tanggal kehadiran tersebut. Lalu list aktivitas karyawan akan dicari berdasarkan id kehadiran. Terdapat query parameter yang bersifat optional pada endpoint ini. Jika tidak ada query parameter yang dimasukkan, maka sistem akan memberikan tanggal default berupa tanggal hari tersebut, sedangkan jika terdapat tanggal yang dimasukkan pada query parameter, maka tanggal tersebut akan dipakai untuk dicari datanya.

#### HTTP Request
```json
GET https://go-attendance.herokuapp.com/activity?activity_date={activity_date}
```

#### Parameters

| Body    |               | Description  |
| ------------- |:-------------:| -------------|
| activity_date   | optional	  	| tanggal dari aktivitas-aktivitas yang ingin dicari dengan format `YYYY-MM-dd`|


#### Result

| Parameters    |  Description  |
| ------------- |:--------------|
|id| `id` dari data aktivitas yang berhasil didapat|
|attendance_id| `id` dari data kehadiran yang berhasil didapat|
|user_id| `id` dari data user yang berhasil didapat|
|description| deskripsi dari aktivitas yang berhasil didapat|
|created_at| tanggal dibuatnya data aktivitas yang berhasil didapat|
|updated_at| tanggal diubah data aktivitas yang berhasil didapat|
|total_activities| total dari data aktivitas yang berhasil didapat|

#### Example
```json
curl https://go-attendance.herokuapp.com/activity -X GET -H "Content-Type: application/json"

or

curl https://go-attendance.herokuapp.com/activity?activity_date=2023-01-17 -X GET -H "Content-Type: application/json"
```
```json
{
    "data": {
        "activities": [
            {
                "id": 1,
                "attendance_id": 2,
                "user_id": 1,
                "description": "Pada hari ini memperbaiki bug yang ada",
                "created_at": "2023-01-16T04:13:52+07:00",
                "updated_at": "2023-01-16T04:13:52+07:00"
            },
            {
                "id": 2,
                "attendance_id": 2,
                "user_id": 1,
                "description": "Pada hari ini, bugnya sudah selesai diperbaiki",
                "created_at": "2023-01-16T04:31:56+07:00",
                "updated_at": "2023-01-16T04:31:56+07:00"
            }
        ],
        "total_activities": 2
    }
}
```