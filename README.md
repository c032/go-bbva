# go-bbva

## Executables

### `bbva-to-json`

Usage:

```
bbva-to-json exported.xlsx
```

That command prints one JSON object per line, for example:

```
{"Concepto":"Example.","Disponible":"0.0","Divisa":"EUR","F.Valor":"12/07/2024","Fecha":"12/07/2024","Importe":"0","Movimiento":"Example.","Observaciones":"Example."}
```

Keys and values are **not** parsed, this only converts an `.xlsx` file
to a simplified JSON that can be used for further processing.

## License

MPL 2.0.
