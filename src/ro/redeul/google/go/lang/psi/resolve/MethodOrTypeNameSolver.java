package ro.redeul.google.go.lang.psi.resolve;

import ro.redeul.google.go.lang.psi.declarations.GoVarDeclaration;
import ro.redeul.google.go.lang.psi.expressions.literals.GoLiteralIdentifier;
import ro.redeul.google.go.lang.psi.processors.ResolveStates;
import ro.redeul.google.go.lang.psi.resolve.references.AbstractCallOrConversionReference;
import ro.redeul.google.go.lang.psi.statements.GoShortVarDeclaration;
import ro.redeul.google.go.lang.psi.toplevel.GoFunctionDeclaration;
import ro.redeul.google.go.lang.psi.toplevel.GoFunctionParameter;
import ro.redeul.google.go.lang.psi.toplevel.GoTypeSpec;
import ro.redeul.google.go.lang.psi.types.GoPsiType;
import ro.redeul.google.go.lang.psi.types.GoPsiTypeFunction;

public class MethodOrTypeNameSolver<
        Ref extends AbstractCallOrConversionReference<Solver, Ref>,
        Solver extends MethodOrTypeNameSolver<Ref, Solver>>
    extends RefSolver<Ref, Solver> {
    public MethodOrTypeNameSolver(Ref reference) {
        super(reference);
    }

    @Override
    public void visitFunctionDeclaration(GoFunctionDeclaration declaration) {
        if (checkReference(declaration))
            addTarget(declaration, declaration.getNameIdentifier());
    }

    @Override
    public void visitTypeSpec(GoTypeSpec type) {

        if (ResolveStates.get(getState(), ResolveStates.Key.IsPackageBuiltin)) {
            String typeName = type.getName();
            GoPsiType typeDeclaration = type.getType();
            if (typeName != null && typeDeclaration != null) {
                if (!typeName.equals(typeDeclaration.getText()))
                    return;
            }
        }

        if (checkReference(type.getTypeNameDeclaration()))
            addTarget(type);
    }

    @Override
    public void visitVarDeclaration(GoVarDeclaration declaration) {
        if (checkReference(declaration))
            addTarget(declaration);
    }

    @Override
    public void visitShortVarDeclaration(GoShortVarDeclaration declaration) {
        GoLiteralIdentifier ids[] = declaration.getDeclarations();
        checkIdentifiers(getReferenceName(), ids);
    }

    private boolean checkVarDeclaration(GoShortVarDeclaration declaration) {
        declaration.getIdentifiersType();
        return false;
    }

    @Override
    public void visitFunctionParameter(GoFunctionParameter parameter) {
        if (!(parameter.getType() instanceof GoPsiTypeFunction)) {
            return;
        }

        for (GoLiteralIdentifier identifier : parameter.getIdentifiers()) {
            if (!checkReference(identifier)) {
                continue;
            }

            if (!addTarget(identifier)) {
                return;
            }
        }
    }
}
