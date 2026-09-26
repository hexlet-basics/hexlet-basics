{{/*
Shared fragments of the codebasics chart. Selector labels are immutable on a
Deployment, so a workload that keeps a legacy name (web) keeps its legacy
selector too.
*/}}

{{/* Top-level labels (metadata.labels) of every resource. */}}
{{- define "codebasics.labels" -}}
app.kubernetes.io/name: "{{ .Chart.Name }}"
app.kubernetes.io/instance: "{{ .Release.Name }}"
app.kubernetes.io/version: "{{ .Chart.Version }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
{{- end -}}

{{/* nodeSelector + imagePullSecrets: every pod runs on the codebasics nodes. */}}
{{- define "codebasics.podScheduling" -}}
nodeSelector:
  group: "{{ .Release.Namespace }}"
imagePullSecrets:
  - name: docker-config
{{- end -}}

{{/* nodeAffinity for a node role (web or service). */}}
{{- define "codebasics.nodeAffinity" -}}
nodeAffinity:
  requiredDuringSchedulingIgnoredDuringExecution:
    nodeSelectorTerms:
      - matchExpressions:
          - key: role
            operator: In
            values:
              - {{ . }}
{{- end -}}

{{/* Image references, tagged with the chart version. */}}
{{- define "codebasics.goImage" -}}
"{{ .Values.images.go.repository }}:{{ .Chart.Version }}"
{{- end -}}
{{- define "codebasics.webImage" -}}
"{{ .Values.images.web.repository }}:{{ .Chart.Version }}"
{{- end -}}

{{/*
Environment of a Go process: the secret via envFrom plus the plain settings.
`env` wins over envFrom on a duplicate key, so values.yaml overrides whatever
the legacy secret still holds under the same name.
*/}}
{{- define "codebasics.goEnv" -}}
envFrom:
  - secretRef:
      name: codebasics-environment-secrets
env:
  - name: HEXLET_BASICS_RELEASE_VERSION
    value: "{{ .Chart.Version }}"
  {{- range $key, $value := .Values.goEnv }}
  - name: {{ $key }}
    value: {{ $value | quote }}
  {{- end }}
{{- end -}}
