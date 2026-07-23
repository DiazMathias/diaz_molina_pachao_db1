package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"strconv"
	"context"
)

const connStr = "user=postgres host=localhost dbname=diaz_molina_pachao_db1 sslmode=disable"
const connStrPostgres = "user=postgres host=localhost dbname=postgres sslmode=disable"

func conectar() *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	for {
		fmt.Println("\n+++ Menu +++")
		fmt.Println("1 Crear base de datos")
		fmt.Println("2 Crear tablas")
		fmt.Println("3 Agregar PKs y FKs")
		fmt.Println("4 Eliminar PKs y FKs")
		fmt.Println("5 Cargar datos")
		fmt.Println("6 Crear stored procedures y triggers")
		fmt.Println("7 Iniciar pruebas")
		fmt.Println("8 Cargar datos en BoltDB")
		fmt.Println("0 Salir")
		fmt.Print("Eliga la opcion: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			crearBaseDeDatos()
		case 2:
			crearTablas()
		case 3:
			establecerPKs()
			establecerFKs()
		case 4:
			eliminarFKs()
			eliminarPKs()
		case 5:
			cargarTarifas()
			cargarClientes()
			cargarDirecciones()
			cargarProductos()
			cargarDatosDePrueba()
		case 6:
			spCrearPedido()
			spAgregarProducto()
			spEntregarPedido()
			spAnularPedido()
			trgEnviarEmailPedidoIniciado()
			trgEnviarEmailPedidoEntregado()
			trgEnviarEmailPedidoAnulado()
		case 7:
			iniciarPruebas()
		case 8:
			cargarDatosBoltDB()
			revisarDatosBoltDB()
		case 0:
			fmt.Println("\nSaliendo...")
			os.Exit(0)
		default:
			fmt.Println("\nOpcion invalida")
		}
	}
}

func crearBaseDeDatos() {
	db, err := sql.Open("postgres", connStrPostgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("create database diaz_molina_pachao_db1")
	if err != nil {
		fmt.Println("\nNo se pudo crear la base de datos", err)
	} else {
		fmt.Println("\nBase de datos creada!")
	}
}

func crearTablas() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/tablas.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo tablas.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al ejecutar tablas.sql: %v\n", err)
		return
	}

	fmt.Println("\nTablas creadas!")
}

func establecerPKs() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/pks.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo pks.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al ejecutar pks.sql: %v\n", err)
		return
	}

	fmt.Println("\nClaves primarias establecidas!")
}

func establecerFKs() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/fks.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo fks.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al ejecutar fks.sql: %v\n", err)
		return
	}

	fmt.Println("\nClaves foràneas establecidas!")
}

func eliminarFKs() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/eliminar_fks.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo eliminar_fks: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al ejecutar eliminar_fks: %v\n", err)
		return
	}

	fmt.Println("\nClaves foràneas eliminadas!")
}

func eliminarPKs() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/eliminar_pks.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo eliminar_pks: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al ejecutar eliminar_pks: %v\n", err)
		return
	}

	fmt.Println("\nClaves primarias eliminadas!")
}

func cargarTarifas() {
	db := conectar()
	defer db.Close()

	bytes, err := os.ReadFile("json/tarifas_entrega.json")
	if err != nil {
		log.Fatal(err)
	}

	type TarifaEntrega struct {
		CodigoPostal string  `json:"codigo_postal"`
		CostoDecimal float64 `json:"costo"`
	}

	var listaTarifasEntrega []TarifaEntrega
	err = json.Unmarshal(bytes, &listaTarifasEntrega)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range listaTarifasEntrega {
		query := "insert into tarifa_entrega values ($1, $2)"
		_, err := db.Exec(query, t.CodigoPostal, t.CostoDecimal)
		if err != nil {
			fmt.Printf("\nNo se cargó la tarifa debido a su código postal (CP: %s) repetido.\n", t.CodigoPostal)
		}
	}

	fmt.Println("\nDatos de las tarifas de entrega cargados!")
}

func cargarClientes() {
	db := conectar()
	defer db.Close()

	bytes, err := os.ReadFile("json/clientes.json")
	if err != nil {
		log.Fatal(err)
	}

	type Cliente struct {
		ID              int    `json:"id_usuarie"`
		Nombre          string `json:"nombre"`
		Apellido        string `json:"apellido"`
		Dni             int    `json:"dni"`
		FechaNacimiento string `json:"fecha_nacimiento"`
		Telefono        string `json:"telefono"`
		Email           string `json:"email"`
	}

	var listaClientes []Cliente
	err = json.Unmarshal(bytes, &listaClientes)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range listaClientes {
		query := "insert into cliente values($1, $2, $3, $4, $5, $6, $7)"
		_, err := db.Exec(query, c.ID, c.Nombre, c.Apellido, c.Dni, c.FechaNacimiento, c.Telefono, c.Email)
		if err != nil {
			fmt.Printf("\nNo se cargó el cliente %s debido a su ID (ID: %d) repetido.\n", c.Nombre, c.ID)
		}
	}

	fmt.Println("\nDatos de los clientes cargados!")
}

func cargarDirecciones() {
	db := conectar()
	defer db.Close()

	bytes, err := os.ReadFile("json/direcciones.json")
	if err != nil {
		log.Fatal(err)
	}

	type DireccionEntrega struct {
		IDUsuario          int    `json:"id_usuarie"`
		IDDireccionEntrega int    `json:"id_direccion_entrega"`
		Direccion          string `json:"direccion"`
		Localidad          string `json:"localidad"`
		CodigoPostal       string `json:"codigo_postal"`
	}

	var listaDirecciones []DireccionEntrega
	err = json.Unmarshal(bytes, &listaDirecciones)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range listaDirecciones {
		query := "insert into direccion_entrega values ($1,$2,$3,$4,$5)"
		_, err := db.Exec(query, d.IDUsuario, d.IDDireccionEntrega, d.Direccion, d.Localidad, d.CodigoPostal)
		if err != nil {
			fmt.Printf("\nNo se cargó la dirección %s porque su código postal (%s) no está registrado en las tarifas de entrega.\n", d.Direccion, d.CodigoPostal)
		}
	}

	fmt.Println("\nDatos de las direcciones cargados!")
}

func cargarProductos() {
	db := conectar()
	defer db.Close()

	bytes, err := os.ReadFile("json/productos.json")
	if err != nil {
		log.Fatal(err)
	}

	type Producto struct {
		ID              int     `json:"id_producto"`
		Nombre          string  `json:"nombre"`
		Precio          float64 `json:"precio_unitario"`
		StockDisponible int     `json:"stock_disponible"`
		StockReservado  int     `json:"stock_reservado"`
	}

	var listaProductos []Producto
	err = json.Unmarshal(bytes, &listaProductos)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range listaProductos {
		query := "insert into producto values ($1,$2,$3,$4,$5)"
		_, err := db.Exec(query, p.ID, p.Nombre, p.Precio, p.StockDisponible, p.StockReservado)
		if err != nil {
			fmt.Printf("\nNo se cargó el producto %s porque su código (%d) està repetido.\n", p.Nombre, p.ID)
		}
	}

	fmt.Println("\nDatos de los productos cargados!")
}

func cargarDatosDePrueba() {
	db := conectar()
	defer db.Close()

	bytes, err := os.ReadFile("json/datos_de_prueba.json")
	if err != nil {
		log.Fatal(err)
	}

	type DatoPrueba struct {
		IDOrden            int     `json:"id_orden"`
		Operacion          string  `json:"operacion"`
		IDUsuarie          *int    `json:"id_usuarie"`
		IDDireccionEntrega *int    `json:"id_direccion_entrega"`
		IDPedido           *int    `json:"id_pedido"`
		IDProducto         *int    `json:"id_producto"`
		Cantidad           *int    `json:"cantidad"`
		FechaHoraEntrega   *string `json:"fecha_hora_entrega"`
	}

	var listaDatosPrueba []DatoPrueba
	err = json.Unmarshal(bytes, &listaDatosPrueba)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range listaDatosPrueba {
		query := `insert into datos_de_prueba 
                  (id_orden, operacion, id_usuarie, id_direccion_entrega, id_pedido, id_producto, cantidad, fecha_hora_entrega) 
                  values ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := db.Exec(query, d.IDOrden, d.Operacion, d.IDUsuarie, d.IDDireccionEntrega, d.IDPedido, d.IDProducto, d.Cantidad, d.FechaHoraEntrega)
		if err != nil {
			fmt.Printf("\nNo se cargó el dato de prueba %d (%s): %v\n", d.IDOrden, d.Operacion, err)
		}
	}

	fmt.Println("\nDatos de prueba cargados correctamente!")
}

func spCrearPedido() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/crear_pedido.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo crear_pedido.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el stored procedure crear_pedido: %v\n", err)
		return
	}

	fmt.Println("\nStored procedure de creaciòn de pedido creado!")
}

func spAgregarProducto() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/agregar_producto.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo agregar_producto.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el stored procedure agregar_producto: %v\n", err)
		return
	}

	fmt.Println("\nStored procedure de agregar producto creado!")
}

func spEntregarPedido() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/entregar_pedido.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo entregar_pedido.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el stored procedure entregar_producto: %v\n", err)
		return
	}

	fmt.Println("\nStored procedure de entregar pedido creado!")
}

func spAnularPedido() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/anular_pedido.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo anular_pedido.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el stored procedure anular_pedido: %v\n", err)
		return
	}

	fmt.Println("\nStored procedure de anulación de pedido creado!")
}

func trgEnviarEmailPedidoIniciado() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/enviar_email_pedido_iniciado.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo enviar_email_pedido_iniciado.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el trigger de pedido iniciado: %v\n", err)
		return
	}

	fmt.Println("\nTrigger de email (pedido iniciado) creado con éxito!")
}

func trgEnviarEmailPedidoEntregado() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/enviar_email_pedido_entregado.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo enviar_email_pedido_iniciado.sol: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el trigger de pedido entregado: %v\n", err)
		return
	}

	fmt.Println("\nTrigger de email (pedido entregado) creado con exito!")
}

func trgEnviarEmailPedidoAnulado() {
	db := conectar()
	defer db.Close()

	sql, err := os.ReadFile("sql/enviar_email_pedido_anulado.sql")
	if err != nil {
		fmt.Printf("\nError al leer el archivo enviar_email_pedido_anulado.sql: %v\n", err)
		return
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		fmt.Printf("\nError al crear el trigger de pedido anulado: %v\n", err)
		return
	}

	fmt.Println("\nTrigger de email (pedido anulado) creado con exito!")
}

func iniciarPruebas() {
	db := conectar()
	defer db.Close()

	query := "select id_orden, trim(operacion), id_usuarie, id_direccion_entrega, id_pedido, id_producto, cantidad, fecha_hora_entrega from datos_de_prueba order by id_orden"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("\nError al consultar datos_de_prueba: %v\n", err)
		return
	}
	defer rows.Close()

	type DatoPrueba struct {
		IDOrden            int
		Operacion          string
		IDUsuarie          *int
		IDDireccionEntrega *int
		IDPedido           *int
		IDProducto         *int
		Cantidad           *int
		FechaHoraEntrega   *string
	}

	var pruebas []DatoPrueba

	for rows.Next() {
		var d DatoPrueba
		err := rows.Scan(
			&d.IDOrden,
			&d.Operacion,
			&d.IDUsuarie,
			&d.IDDireccionEntrega,
			&d.IDPedido,
			&d.IDProducto,
			&d.Cantidad,
			&d.FechaHoraEntrega,
		)
		if err != nil {
			fmt.Printf("\nError al escanear fila: %v\n", err)
			continue
		}
		pruebas = append(pruebas, d)
	}

	for _, p := range pruebas {
		if p.Operacion == "creación" {
			tx, err := db.BeginTx(context.Background(), nil)
			if err != nil {
				fmt.Printf("\nError al iniciar transacción para la orden %d: %v\n", p.IDOrden, err)
				continue
			}
			querySP := "select crear_pedido($1::int, $2::int)"
			_, err = tx.Exec(querySP, p.IDUsuarie, p.IDDireccionEntrega)
			if err != nil {
				tx.Rollback()
				fmt.Printf("\nError al ejecutar crear_pedido para la orden %d: %v\n", p.IDOrden, err)
			} else {
				tx.Commit()
				fmt.Printf("\nPedido creado exitosamente desde la orden %d.\n", p.IDOrden)
			}
		}
	}

	for _, p := range pruebas {
		if p.Operacion == "producto" {
			tx, err := db.BeginTx(context.Background(), nil)
			if err != nil {
				fmt.Printf("\nError al iniciar transacción para la orden %d: %v\n", p.IDOrden, err)
				continue
			}
			querySP := "select agregar_producto($1::int, $2::int, $3::int)"
			_, err = tx.Exec(querySP, p.IDPedido, p.IDProducto, p.Cantidad)
			if err != nil {
				tx.Rollback()
				fmt.Printf("\nError al ejecutar agregar_producto para la orden %d: %v\n", p.IDOrden, err)
			} else {
				tx.Commit()
				fmt.Printf("\nProducto agregado exitosamente desde la orden %d.\n", p.IDOrden)
			}
		}
	}

	for _, p := range pruebas {
		if p.Operacion == "entrega" {
			tx, err := db.BeginTx(context.Background(), nil)
			if err != nil {
				fmt.Printf("\nError al iniciar transacción para la orden %d: %v\n", p.IDOrden, err)
				continue
			}
			querySP := "select entregar_pedido($1::int, $2::timestamp)"
			_, err = tx.Exec(querySP, p.IDPedido, p.FechaHoraEntrega)
			if err != nil {
				tx.Rollback()
				fmt.Printf("\nError al ejecutar entregar_pedido para la orden %d: %v\n", p.IDOrden, err)
			} else {
				tx.Commit()
				fmt.Printf("\nPedido entregado exitosamente desde la orden %d.\n", p.IDOrden)
			}
		}

		if p.Operacion == "anulación" {
			tx, err := db.BeginTx(context.Background(), nil)
			if err != nil {
				fmt.Printf("\nError al iniciar transacción para la orden %d: %v\n", p.IDOrden, err)
				continue
			}
			querySP := "select anular_pedido($1::int)"
			_, err = tx.Exec(querySP, p.IDPedido)
			if err != nil {
				tx.Rollback()
				fmt.Printf("\nError al ejecutar anular_pedido para la orden %d: %v\n", p.IDOrden, err)
			} else {
				tx.Commit()
				fmt.Printf("\nPedido anulado exitosamente desde la orden %d.\n", p.IDOrden)
			}
		}
	}
}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	err = b.Put(key, val)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type ClienteBolt struct {
	ID              int    `json:"id_usuarie"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	DNI             int    `json:"dni"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
}

type DireccionBolt struct {
	IDUsuarie          int    `json:"id_usuarie"`
	IDDireccionEntrega int    `json:"id_direccion_entrega"`
	Direccion          string `json:"direccion"`
	Localidad          string `json:"localidad"`
	CodigoPostal       string `json:"codigo_postal"`
}

type ProductoBolt struct {
	ID              int     `json:"id_producto"`
	Nombre          string  `json:"nombre"`
	PrecioUnitario  float64 `json:"precio_unitario"`
	StockDisponible int     `json:"stock_disponible"`
	StockReservado  int     `json:"stock_reservado"`
}

type DetallePedidoBolt struct {
	IDProducto     int     `json:"id_producto"`
	NombreProducto string  `json:"nombre_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

type PedidoBolt struct {
	IDPedido           int                 `json:"id_pedido"`
	FechaPedido        string              `json:"f_pedido"`
	FechaEntrega       *string             `json:"fecha_entrega"`
	HoraEntrega        *string             `json:"hora_entrega"`
	IDUsuarie          int                 `json:"id_usuarie"`
	IDDireccionEntrega int                 `json:"id_direccion_entrega"`
	MontoTotal         float64             `json:"monto_total"`
	CostoEnvio         float64             `json:"costo_envio"`
	Estado             string              `json:"estado"`
	Detalle            []DetallePedidoBolt `json:"detalle"`
}

func cargarDatosBoltDB() {
	dbBolt, err := bolt.Open("nosql.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbBolt.Close()

	dbSQL := conectar()
	defer dbSQL.Close()

	migrarClientes(dbSQL, dbBolt)
	migrarProductos(dbSQL, dbBolt)
	migrarDirecciones(dbSQL, dbBolt)
	migrarPedidos(dbSQL, dbBolt)
}

func migrarClientes(dbSQL *sql.DB, dbBolt *bolt.DB) {
	filaClientes, err := dbSQL.Query("select id_usuarie, nombre, apellido, dni, fecha_nacimiento, telefono, email from cliente")
	if err != nil {
		fmt.Printf("Error al buscar datos de clientes para BoltDB: %v\n", err)
		return
	}
	defer filaClientes.Close()

	for filaClientes.Next() {
		var c ClienteBolt
		filaClientes.Scan(&c.ID, &c.Nombre, &c.Apellido, &c.DNI, &c.FechaNacimiento, &c.Telefono, &c.Email)

		bytesJSON, _ := json.Marshal(c)
		key := []byte(strconv.Itoa(c.ID))
		CreateUpdate(dbBolt, "clientes", key, bytesJSON)
	}
	fmt.Println("Clientes migrados exitosamente.")
}

func migrarProductos(dbSQL *sql.DB, dbBolt *bolt.DB) {
	filaProductos, err := dbSQL.Query("select id_producto, nombre, precio_unitario, stock_disponible, stock_reservado from producto")
	if err != nil {
		fmt.Printf("Error al buscar datos de productos para BoltDB: %v\n", err)
		return
	}
	defer filaProductos.Close()

	for filaProductos.Next() {
		var p ProductoBolt
		filaProductos.Scan(&p.ID, &p.Nombre, &p.PrecioUnitario, &p.StockDisponible, &p.StockReservado)

		bytesJSON, _ := json.Marshal(p)
		key := []byte(strconv.Itoa(p.ID))
		CreateUpdate(dbBolt, "productos", key, bytesJSON)
	}
	fmt.Println("Productos migrados exitosamente.")
}

func migrarDirecciones(dbSQL *sql.DB, dbBolt *bolt.DB) {
	filaDirecciones, err := dbSQL.Query("select id_usuarie, id_direccion_entrega, direccion, localidad, codigo_postal from direccion_entrega")
	if err != nil {
		fmt.Printf("Error al buscar datos de direcciones para BoltDB: %v\n", err)
		return
	}
	defer filaDirecciones.Close()

	for filaDirecciones.Next() {
		var d DireccionBolt
		filaDirecciones.Scan(&d.IDUsuarie, &d.IDDireccionEntrega, &d.Direccion, &d.Localidad, &d.CodigoPostal)

		bytesJSON, _ := json.Marshal(d)
		key := []byte(strconv.Itoa(d.IDUsuarie) + "-" + strconv.Itoa(d.IDDireccionEntrega))
		CreateUpdate(dbBolt, "direcciones", key, bytesJSON)
	}
	fmt.Println("Direcciones migradas exitosamente.")
}

func migrarPedidos(dbSQL *sql.DB, dbBolt *bolt.DB) {
	filaPedidos, err := dbSQL.Query(`
		select id_pedido, f_pedido, fecha_entrega, hora_entrega,
		       id_usuarie, id_direccion_entrega, monto_total, costo_envio, trim(estado)
		from pedido`)
	if err != nil {
		fmt.Printf("Error al buscar datos de pedidos para BoltDB: %v\n", err)
		return
	}
	defer filaPedidos.Close()

	for filaPedidos.Next() {
		var p PedidoBolt
		filaPedidos.Scan(
			&p.IDPedido, &p.FechaPedido, &p.FechaEntrega, &p.HoraEntrega,
			&p.IDUsuarie, &p.IDDireccionEntrega, &p.MontoTotal, &p.CostoEnvio, &p.Estado,
		)

		filaDetalle, err := dbSQL.Query(`
			select pd.id_producto, pr.nombre, pd.cantidad, pd.precio_unitario
			from pedido_detalle pd
			join producto pr on pd.id_producto = pr.id_producto
			where pd.id_pedido = $1`, p.IDPedido)
		if err != nil {
			fmt.Printf("Error en detalle del pedido %d: %v\n", p.IDPedido, err)
			continue
		}

		for filaDetalle.Next() {
			var det DetallePedidoBolt
			filaDetalle.Scan(&det.IDProducto, &det.NombreProducto, &det.Cantidad, &det.PrecioUnitario)
			p.Detalle = append(p.Detalle, det)
		}
		filaDetalle.Close()

		if len(p.Detalle) < 3 {
			fmt.Printf("Pedido ID %d omitido: Solo tiene %d producto(s) (el minimo requerido es 3).\n", p.IDPedido, len(p.Detalle))
			continue
		}

		bytesJSON, _ := json.Marshal(p)
		key := []byte(strconv.Itoa(p.IDPedido))
		CreateUpdate(dbBolt, "pedidos", key, bytesJSON)
	}
	fmt.Println("Pedidos migrados exitosamente.")
}

func revisarDatosBoltDB() {
	dbBolt, err := bolt.Open("nosql.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbBolt.Close()

	fmt.Println("\nREVISION DATOS BOLTDB")

	buckets := []string{"clientes", "productos", "direcciones", "pedidos"}

	_ = dbBolt.View(func(tx *bolt.Tx) error {
		for _, bName := range buckets {
			b := tx.Bucket([]byte(bName))
			if b == nil {
				fmt.Printf("\nEl bucket '%s' no existe.\n", bName)
				continue
			}

			fmt.Printf("\n--- BUCKET: %s ---\n", bName)

			_ = b.ForEach(func(k, v []byte) error {
				fmt.Printf("Clave: %s\n", string(k))

				var objeto interface{}
				if err := json.Unmarshal(v, &objeto); err == nil {
					bytesFormateados, _ := json.MarshalIndent(objeto, "", "    ")
					fmt.Printf("Valor:\n%s\n", string(bytesFormateados))
				} else {
					fmt.Println("El dato existe en BoltDB pero esta corrupto o no es JSON.")
					fmt.Printf("Valor (Raw): %s\n", string(v))
				}
				fmt.Println("-----------------------------------------")
				return nil
			})
		}
		return nil
	})
}
