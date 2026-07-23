create or replace function trg_enviar_email_pedido_entregado() returns trigger as $$
declare
	var_id_email int;
	var_email_cliente text;
	var_cuerpo_email text;
	var_fila record;
begin
	if NEW.estado = 'entregado' and OLD.estado != 'entregado' then
		select email into var_email_cliente from cliente where id_usuarie = new.id_usuarie;

		var_cuerpo_email := 'El pedido ' || NEW.id_pedido || ' ha sido entregado con exito. Detalle: ';

		select coalesce(max(id_email), 0) + 1 into var_id_email from envio_email;

		for var_fila in
			(select id_producto,
			 cantidad,
			 precio_unitario
			 from pedido_detalle where id_pedido = NEW.id_pedido)
		loop
			var_cuerpo_email := var_cuerpo_email || 'Producto' || var_fila.id_producto ||
								', Cantidad: ' || var_fila.cantidad ||
								', Precio: $' || var_fila.precio_unitario || ' | ';
		end loop;
		
		insert into envio_email (id_email, f_generacion, email_cliente, asunto, cuerpo, estado)
		values (var_id_email, now(), var_email_cliente, 'Pedido entregado', var_cuerpo_email, 'enviado');
		
	end if;
	return new;
end;
$$ language plpgsql;

create or replace trigger trg_pedido_entregado 
after update on pedido
for each row
execute function trg_enviar_email_pedido_entregado();
