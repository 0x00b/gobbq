// NOTE:!!
// DON'T EDIT IT!!
//  请勿添加其他内容，包括函数，结构体，变量等等，否则重新生成时会丢失。

import { UnaryResponse } from "gobbq-ts/dist/src/context/unary";
import { Client } from "gobbq-ts/dist/src";
import { makeClientConstructor } from "gobbq-ts/dist/src/bbq/bbq";
import { ServiceType } from "gobbq-ts/dist/proto/bbq";
import Long from "long";

{{- range $sidx, $m := $.Messages }}
import { {{$m.GoIdent.GoName}} } from "./{{FileName $.Name}}"
{{- end -}}

{{range $sidx, $s := $.Services}}
{{- $sName := $s.GoName -}}
{{- $isSvc := isService $s -}}
{{- $typeName := concat "" $sName "Entity" -}}
{{- if $isSvc}}
	{{$typeName = concat "" $sName "Service"}}
{{end -}}
// {{goComments $typeName $s.Comments}}
export type {{$typeName}}Definition = typeof {{$typeName}}Definition;
export const {{$typeName}}Definition = {
  typeName: "{{$.GoPackageName}}.{{$typeName}}",
  serviceType: ServiceType.{{- if $isSvc}}Service{{else}}Entity{{end -}}, 
  methods: {
	{{- range $midx, $m := $s.Methods }}
    {{$m.GoName}}: {
      methodName: "{{$m.GoName}}",
      requestType: {{$m.GoInput.GoIdent.GoName}},
      responseType: {{if $m.HasResponse}}{{$m.GoOutput.GoIdent.GoName}}{{else}}undefined{{end}},
      requestSerialize: (req: {{$m.GoInput.GoIdent.GoName}}): Buffer => {
        return Buffer.from({{$m.GoInput.GoIdent.GoName}}.encode(req).finish())
      },
      requestDeserialize: (input: Uint8Array): {{$m.GoInput.GoIdent.GoName}} => {
        return {{$m.GoInput.GoIdent.GoName}}.decode(input)
      },
      responseSerialize: (req: {{if $m.HasResponse}}{{$m.GoOutput.GoIdent.GoName}}{{else}}any{{end}}): Buffer => {
        return {{if $m.HasResponse}}Buffer.from({{$m.GoOutput.GoIdent.GoName}}.encode(req).finish()){{else}}Buffer.from(""){{end}}
      },
      responseDeserialize: (input: Uint8Array): {{if $m.HasResponse}}{{$m.GoOutput.GoIdent.GoName}}{{else}}any{{end}} => {
        {{if $m.HasResponse}}return {{$m.GoOutput.GoIdent.GoName}}.decode(input){{end}}
      },
    },
	{{end -}}
  },
} as const;

export interface {{$typeName}} {
{{range $midx, $m := $s.Methods}}
	// {{goComments $m.GoName $m.Comments}}
	{{$m.GoName}}(request: {{$m.GoInput.GoIdent.GoName}}){{if $m.HasResponse}}:UnaryResponse<{{$m.GoOutput.GoIdent.GoName}}>{{else}}:void{{end}}
{{end -}}
}

export function New{{$typeName}}(client: Client<any> {{- if $isSvc}}{{else}}, entityID: Long{{end -}}): {{$typeName}} {
  return makeClientConstructor(client, {{$typeName}}Definition{{- if $isSvc}}{{else}}, entityID{{end -}}) as unknown as {{$typeName}}
}
{{end -}}
