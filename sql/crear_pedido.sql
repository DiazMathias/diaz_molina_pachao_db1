create or replace function crear_pedido(p_id_usuarie int, p_id_direccion int) returns boolean as $$
declare
    var_cp char(4);
    var_costo_envio decimal(12,2);
    var_id_pedido int;
    var_id_error int;
begin
    if not exists (select 1 from cliente where id_usuarie = p_id_usuarie) then
        select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, operacion, motivo)
        values (var_id_error, 'creaciòn', 'id de usuarie no válido');
        return false;
    end if;

    if not exists (select 1 from direccion_entrega where id_usuarie = p_id_usuarie and id_direccion_entrega = p_id_direccion) then
        select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, operacion, motivo)
        values (var_id_error, 'creaciòn', 'id de dirección no válido');
        return false;
    end if;

    select codigo_postal into var_cp from direccion_entrega where id_usuarie = p_id_usuarie and id_direccion_entrega = p_id_direccion;

    if not exists (select 1 from tarifa_entrega where codigo_postal = var_cp) then
        select coalesce(max(id_error), 0) + 1 into var_id_error from error;
        insert into error (id_error, operacion, motivo)
        values (var_id_error, 'creaciòn', 'dirección de entrega fuera del área de atención');
        return false;
    end if;

    select costo into var_costo_envio from tarifa_entrega where codigo_postal = var_cp;

    select coalesce(max(id_pedido), 0) + 1 into var_id_pedido from pedido;

    insert into pedido (id_pedido, id_usuarie, id_direccion_entrega, f_pedido, costo_envio, estado)
    values (var_id_pedido, p_id_usuarie, p_id_direccion, now(), var_costo_envio, 'ingresado');
    return true;
end;
$$ language plpgsql;
