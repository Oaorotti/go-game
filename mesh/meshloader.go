package mesh

import (
	"github.com/bloeys/assimp-go/asig"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Mesh struct {
	VAOs       []uint32
	NumIndices []int
	Release    func()
}

func createBuffersAndVAO(mesh *asig.Mesh) (vao, vbo, ebo uint32) {
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	vertices := make([]float32, 0)
	indices := make([]uint32, 0)

	for _, vertex := range mesh.Vertices {
		vertices = append(vertices, vertex.X(), vertex.Y(), vertex.Z())
	}
	for _, face := range mesh.Faces {
		indices = append(indices, uint32(face.Indices[0]), uint32(face.Indices[1]), uint32(face.Indices[2]))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	return vao, vbo, ebo
}

func LoadMesh(path string) *Mesh {
	scene, release, err := asig.ImportFile(path, asig.PostProcessTriangulate|asig.PostProcessJoinIdenticalVertices)

	if err != nil {
		panic(err)
	}

	var meshVAOs []uint32
	var meshNumIndices []int

	for _, mesh := range scene.Meshes {
		vao, _, _ := createBuffersAndVAO(mesh)
		meshVAOs = append(meshVAOs, vao)
		meshNumIndices = append(meshNumIndices, len(mesh.Faces)*3)
	}

	return &Mesh{
		VAOs:       meshVAOs,
		NumIndices: meshNumIndices,
		Release:    release,
	}
}

func (m *Mesh) ReleaseMesh() {
	m.Release()
}

func RenderMesh(mesh Mesh) {
	for i, vao := range mesh.VAOs {
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, int32(mesh.NumIndices[i]), gl.UNSIGNED_INT, nil)
	}
}
