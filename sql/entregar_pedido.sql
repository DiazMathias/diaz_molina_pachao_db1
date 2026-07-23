create or replace function entregar_pedido(p_id_pedido int, p_fecha_hora_entrega timestamp) returns boolean as $$
declare
	var_estado_actual char(10);
	var_f_pedido timestamp;
	var_id_error int;
begin
    if not exists (select 1 from pedido where id_pedido = p_id_pedido) then
        select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, id_pedido, operacion, motivo)
        values (var_id_error, p_id_pedido, 'entrega', '?id de pedido no valido');
        return false;
    end if;
    
    select trim(estado) into var_estado_actual from pedido where id_pedido = p_id_pedido;
    
	if var_estado_actual != 'ingresado' then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, id_pedido, operacion, motivo)
        values (var_id_error, p_id_pedido, 'entrega', '?pedido ya entregado o anulado');
        return false;
    end if;
    
    if not exists (select 1 from pedido_detalle where id_pedido = p_id_pedido and cantidad > 0) then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, id_pedido, operacion, motivo)
        values (var_id_error, p_id_pedido, 'entrega', '?pedido vacio');
        return false;
    end if;
    
    select f_pedido into var_f_pedido from pedido where id_pedido = p_id_pedido;
    if (p_fecha_hora_entrega <= var_f_pedido) then
		select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, id_pedido, operacion, motivo)
        values (var_id_error, 'entrega', '?fecha de entrega no valida');
        return false;
    end if;
    
    update pedido
    set estado = 'entregado',
		fecha_entrega = p_fecha_hora_entrega::date,
		hora_entrega = p_fecha_hora_entrega::time
	where id_pedido = p_id_pedido;
	
	update producto p	
	set stock_reservado = p.stock_reservado - pd.cantidad
	from pedido_detalle pd
	where pd.id_pedido = p_id_pedido and p.id_producto = pd.id_producto;
	
	return true;
end;
$$ language plpgsql;
	
    
    
