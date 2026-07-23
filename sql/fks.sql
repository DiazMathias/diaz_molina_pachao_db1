alter table direccion_entrega add constraint id_usuarie_fk foreign key (id_usuarie) references cliente (id_usuarie);
alter table direccion_entrega add constraint codigo_postal_fk foreign key (codigo_postal) references tarifa_entrega (codigo_postal);
alter table pedido add constraint id_usuarie_fk foreign key (id_usuarie) references cliente (id_usuarie);
alter table pedido add constraint direccion_entrega_fk foreign key (id_usuarie, id_direccion_entrega) references direccion_entrega (id_usuarie, id_direccion_entrega);
alter table pedido_detalle add constraint id_pedido_fk foreign key (id_pedido) references pedido (id_pedido);
alter table pedido_detalle add constraint id_producto_fk foreign key (id_producto) references producto (id_producto);
