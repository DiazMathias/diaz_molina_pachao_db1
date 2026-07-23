alter table direccion_entrega drop constraint id_usuarie_fk;
alter table direccion_entrega drop constraint codigo_postal_fk;
alter table pedido drop constraint id_usuarie_fk;
alter table pedido drop constraint direccion_entrega_fk;
alter table pedido_detalle drop constraint id_pedido_fk;
alter table pedido_detalle drop constraint id_producto_fk;
