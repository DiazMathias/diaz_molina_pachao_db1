create or replace function anular_pedido(p_id_pedido int) returns boolean as $$
declare
	var_estado_actual char(10);
	var_id_error int;
begin
	if not exists (select 1 from pedido where id_pedido = p_id_pedido) then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, operacion, motivo)
		values (var_id_error, p_id_pedido, 'anulación', '?id de pedido no valido');
		return false;
	end if;

	select trim(estado) into var_estado_actual from pedido where id_pedido = p_id_pedido;

	if var_estado_actual != 'ingresado' then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
		insert into error (id_error, id_pedido, operacion, motivo)
		values (var_id_error, p_id_pedido, 'anulación', '?pedido ya entregado o anulado');
		return false;
	end if;

	-- Sumamos al disponible y restamos del reservado
	update producto p
	set stock_disponible = p.stock_disponible + pd.cantidad,
		stock_reservado  = p.stock_reservado  - pd.cantidad
	from pedido_detalle pd
	where pd.id_pedido = p_id_pedido
	  and p.id_producto = pd.id_producto;

	update pedido
	set estado = 'anulado'
	where id_pedido = p_id_pedido;

	return true;
	
end;
$$ language plpgsql;
