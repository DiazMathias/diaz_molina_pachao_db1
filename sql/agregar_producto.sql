create or replace function agregar_producto(p_id_pedido int, p_id_producto int, p_cantidad int) returns boolean as $$
declare
	var_stock_disponible int;
	var_precio_producto decimal(12,2);
	var_estado_pedido text;
	var_id_error int;
begin
	if not exists (select 1 from pedido where id_pedido = p_id_pedido) then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, id_producto, cantidad, operacion, motivo)
		values (var_id_error, p_id_pedido, p_id_producto, p_cantidad, 'agregado', 'ID de pedido no valido');
		return false;
	end if;

	select trim(estado) into var_estado_pedido from pedido where id_pedido = p_id_pedido;
	if var_estado_pedido != 'ingresado' then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, id_producto, cantidad, operacion, motivo)
		values (var_id_error, p_id_pedido, p_id_producto, p_cantidad, 'agregado', 'pedido cerrado');
		return false;
	end if;

	if not exists (select 1 from producto where id_producto = p_id_producto) then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, id_producto, cantidad, operacion, motivo)
		values (var_id_error, p_id_pedido, p_id_producto, p_cantidad, 'agregado','ID de producto no valido');
		return false;
	end if;

	-- for update para evitar race condition en el stock
	select stock_disponible, precio_unitario into var_stock_disponible, var_precio_producto from producto where id_producto = p_id_producto for update;
	if var_stock_disponible < p_cantidad then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, id_producto, cantidad, operacion, motivo)
		values (var_id_error, p_id_pedido, p_id_producto, p_cantidad, 'agregado', 'stock no disponible para el producto ' || p_id_producto);
		return false;
	end if;

	update producto
	set stock_disponible = stock_disponible - p_cantidad,
	stock_reservado = stock_reservado + p_cantidad
	where id_producto = p_id_producto;
	
	if exists (select 1 from pedido_detalle where id_pedido = p_id_pedido and id_producto = p_id_producto) then
		update pedido_detalle
		set cantidad = cantidad + p_cantidad
		where id_pedido = p_id_pedido and id_producto = p_id_producto;
	else
		insert into pedido_detalle (id_pedido, id_producto, cantidad, precio_unitario)
		values (p_id_pedido, p_id_producto, p_cantidad, var_precio_producto);
	end if;

	update pedido
	set monto_total = (select sum(cantidad * precio_unitario) from pedido_detalle where id_pedido = p_id_pedido)
	where id_pedido = p_id_pedido;
	return true;
end;
$$ language plpgsql;
