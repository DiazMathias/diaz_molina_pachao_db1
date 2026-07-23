create or replace function trg_enviar_email_pedido_iniciado() returns trigger as $$
declare
	var_id_email int;
	var_email_cliente text;
	var_nombre_cliente text;
	var_cuerpo_email text;
begin
	select email, nombre into var_email_cliente, var_nombre_cliente from cliente where id_usuarie = new.id_usuarie;

	var_cuerpo_email := 'Hola ' || var_nombre_cliente || ', tu pedido con id ' || new.id_pedido || 
	                    ' ha sido iniciado con éxito en la fecha ' || new.f_pedido || 
	                    '. El costo de envío asignado es de $' || new.costo_envio || '.';

	select coalesce(max(id_email), 0) + 1 into var_id_email from envio_email;

	insert into envio_email (id_email, f_generacion, email_cliente, asunto, cuerpo, estado)
	values (var_id_email, now(), var_email_cliente, 'Pedido iniciado', var_cuerpo_email, 'pendiente');

	return new;
end;
$$ language plpgsql;

create or replace trigger trg_pedido_iniciado 
after insert on pedido
for each row
execute function trg_enviar_email_pedido_iniciado();
