create or replace function trg_enviar_email_pedido_anulado() returns trigger as $$
declare
	var_id_email      int;
	var_email_cliente text;
	var_nombre_cliente text;
	var_cuerpo_email  text;
	var_fila          record;
begin
	if NEW.estado = 'anulado' and OLD.estado != 'anulado' then
		select email, nombre into var_email_cliente, var_nombre_cliente
		from cliente
		where id_usuarie = NEW.id_usuarie;

		var_cuerpo_email := 'Hola ' || var_nombre_cliente ||
		                    ', tu pedido ' || NEW.id_pedido ||
		                    ' ha sido anulado. Detalle de productos: ';

		for var_fila in (
			select id_producto, cantidad, precio_unitario
			from pedido_detalle
			where id_pedido = NEW.id_pedido
		) loop
			var_cuerpo_email := var_cuerpo_email ||
			                    'Producto ' || var_fila.id_producto ||
			                    ', Cantidad: ' || var_fila.cantidad ||
			                    ', Precio: $' || var_fila.precio_unitario || ' | ';
		end loop;

		select coalesce(max(id_email), 0) + 1 into var_id_email from envio_email;

		insert into envio_email (id_email, f_generacion, email_cliente, asunto, cuerpo, estado)
		values (var_id_email, now(), var_email_cliente, 'Pedido anulado', var_cuerpo_email, 'enviado');
	end if;

	return new;
end;
$$ language plpgsql;

create or replace trigger trg_pedido_anulado
after update on pedido
for each row
execute function trg_enviar_email_pedido_anulado();
